import binascii
import os
import signal
import six
import subprocess
import time


ROOT_DIR = os.path.dirname(os.path.dirname(__file__))

META_FILE_CONTENT = """package buildmeta

const (
	Version     = "{version}"
	GitCommitID = "{git_hash}"
)
"""
SCHEMA_FILE_CONTENT = """package docs

var doc = "" +
	"{string}"
"""


def gen_next_version(file_path):
	with open(file_path, "r") as version_file:
		cur_version = version_file.readline().strip()
		version_tuple = cur_version.split(".")
		if len(version_tuple) != 3:
			raise Exception("Version file must be in format `x.y.z` (not `%s`)" % cur_version)

		version_tuple[-1] = str(int(version_tuple[-1]) + 1)

	new_version = ".".join(version_tuple)
	with open(file_path, "w") as version_file:
		version_file.write(new_version)

	return new_version


def run_cmd(cmd, shell=True, text=True, stderr=subprocess.STDOUT, **kwargs):
	try:
		data = subprocess.check_output(
			cmd,
			shell=shell, text=text, stderr=stderr,
			**kwargs,
		)
		exitcode = 0
	except subprocess.CalledProcessError as exc:
		data = exc.output
		exitcode = exc.returncode
	if data[-1:] == '\n':
		data = data[:-1]
	return exitcode, data


def sync_swagger_files():
	for root, dirs, _ in os.walk(os.path.join(ROOT_DIR, "api", "services")):
		for service_folder in dirs:
			swagger_path = os.path.join(root, service_folder, "docs", "swagger.yaml")
			if not os.path.exists(swagger_path):
				continue

			schema_path = os.path.join(root, service_folder, "docs", "schema.go")
			with open(swagger_path, "r") as swagger_file:
				with open(schema_path, "w") as schema_file:
					swagger_content = swagger_file.read()
					lines_content = '\\n" +\n\t"'.join(swagger_content.replace('"', '\\"').split("\n"))
					schema_file.write(SCHEMA_FILE_CONTENT.format(string=lines_content))
			print("Sync Swagger schema: %s" % swagger_path)


def build(name, filename, env):
	env.update(os.environ)
	print("Building for %s..." % name)

	start_time = time.time()
	cmd = 'go build -ldflags "-s -w" -o %s .' % os.path.join(ROOT_DIR, "bin", filename)
	code, output = run_cmd(cmd, env=env)
	assert code == 0, "Build %s error %s: \n%s" % (name, code, output)

	latency = time.time() - start_time
	print("Building for %s finished in %.2fs" % (name, latency))


def restore_meta_file():
	with open(os.path.join(ROOT_DIR, "buildmeta", "meta.go"), "w") as meta_file:
		meta_file.write(META_FILE_CONTENT.format(
			version="0.0.0",
			git_hash="<unknown>",
		))


def cancel_signal_handler(sig, frame):
	restore_meta_file()


if __name__ == "__main__":
	base_dir = os.path.join(os.path.dirname(__file__), "..")
	os.chdir(base_dir)

	code, git_commit_id = run_cmd("git rev-parse HEAD")
	assert code == 0, "cannot fetch Git commit hash"

	version = gen_next_version(os.path.join(ROOT_DIR, "bin", "version.txt"))
	meta_data = {
		"version": version,
		"git_hash": git_commit_id
	}
	print("Build version: %s" % version)

	signal.signal(signal.SIGINT, cancel_signal_handler)
	signal.signal(signal.SIGTERM, cancel_signal_handler)

	with open(os.path.join(ROOT_DIR, "buildmeta", "meta.go"), "w") as meta_file:
		meta_file.write(META_FILE_CONTENT.format(**meta_data))

	try:
		sync_swagger_files()

		linux_env = {
			"GOOS": "linux",
			"GOARCH": "amd64",
		}
		build("Linux", "torque-go.bin", linux_env)

		# windows_env = {
		#     "GOOS": "windows",
		#     "GOARCH": "amd64",
		# }
		# build("Windows", "torque-go.exe", windows_env)

		print("Build finished.")
	finally:
		restore_meta_file()
