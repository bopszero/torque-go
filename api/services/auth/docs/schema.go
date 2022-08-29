package docs

var doc = "" +
	"openapi: 3.0.0\n" +
	"info:\n" +
	"  version: 1.0.0\n" +
	"  title: Torque Auth Service\n" +
	"servers:\n" +
	"  - url: \"/\"\n" +
	"components:\n" +
	"  securitySchemes:\n" +
	"    CommitAuth:\n" +
	"      type: http\n" +
	"      scheme: bearer\n" +
	"    AccessAuth:\n" +
	"      type: http\n" +
	"      scheme: bearer\n" +
	"    RefreshAuth:\n" +
	"      type: http\n" +
	"      scheme: bearer\n" +
	"  schemas:\n" +
	"    Errors:\n" +
	"      type: array\n" +
	"      items:\n" +
	"        type: string\n" +
	"        example: Field validation for 'type' failed on the 'required' tag.\n" +
	"    ResponseMessage:\n" +
	"      type: string\n" +
	"      example: \"Field validation for 'user_id' failed on the 'required' tag.\"\n" +
	"    ID:\n" +
	"      type: integer\n" +
	"      minimum: 1\n" +
	"      example: 8466032601\n" +
	"    Username:\n" +
	"      type: string\n" +
	"      example: powerman\n" +
	"    AuthCode:\n" +
	"      type: string\n" +
	"      example: \"123456\"\n" +
	"      minLength: 6\n" +
	"      maxLength: 6\n" +
	"    DeviceUID:\n" +
	"      type: string\n" +
	"      example: \"android-0123abcd\"\n" +
	"      maxLength: 128\n" +
	"    JwtToken:\n" +
	"      type: string\n" +
	"      example: \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c\"\n" +
	"    Email:\n" +
	"      type: string\n" +
	"      example: torquebot1@gmail.com\n" +
	"    CountryCode:\n" +
	"      type: string\n" +
	"      example: VN\n" +
	"    Status:\n" +
	"      type: integer\n" +
	"      enum:\n" +
	"        - 1\n" +
	"        - 2\n" +
	"        - 5\n" +
	"        - 3\n" +
	"        - 6\n" +
	"        - 7\n" +
	"      description: |\n" +
	"        * 1 - Init\n" +
	"        * 2 - Pending Analysis\n" +
	"        * 5 - Pending Approval\n" +
	"        * 3 - Approved\n" +
	"        * 6 - Rejected\n" +
	"        * 7 - Failed\n" +
	"    Timestamp:\n" +
	"      type: integer\n" +
	"      example: 1596517386\n" +
	"    Note:\n" +
	"      type: string\n" +
	"      example: Test ...\n" +
	"    Date:\n" +
	"      type: string\n" +
	"      example: \"2006-01-02\"\n" +
	"    UserType:\n" +
	"      type: integer\n" +
	"      enum:\n" +
	"        - 1\n" +
	"        - 2\n" +
	"        - 3\n" +
	"      description: |\n" +
	"        * 1 - Old\n" +
	"        * 2 - Middle\n" +
	"        * 3 - New\n" +
	"    Code:\n" +
	"      type: string\n" +
	"      example: 65692274f18911eaa53dfc3fdb8a2f0e\n" +
	"    Nationality:\n" +
	"      type: string\n" +
	"      example: VN\n" +
	"    IP:\n" +
	"      type: string\n" +
	"      example: 115.73.221.170\n" +
	"    Country:\n" +
	"      type: string\n" +
	"      example: Viet Nam\n" +
	"    Number:\n" +
	"      type: integer\n" +
	"      example: 1\n" +
	"security:\n" +
	"  - AccessAuth: []\n" +
	"tags:\n" +
	"  - name: KYC\n" +
	"  - name: S2S\n" +
	"  - name: Meta\n" +
	"paths:\n" +
	"  /v1/login/input/prepare/:\n" +
	"    post:\n" +
	"      summary: Prepare data for a secure login\n" +
	"      tags:\n" +
	"        - Login\n" +
	"      security: []\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                username:\n" +
	"                  $ref: \"#/components/schemas/Username\"\n" +
	"              required:\n" +
	"                - username\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      nonce_id:\n" +
	"                        $ref: \"#/components/schemas/Code\"\n" +
	"                      nonce:\n" +
	"                        type: string\n" +
	"                        example: e00a37f05518e155f016d327\n" +
	"                      salt:\n" +
	"                        type: string\n" +
	"                        example: 88497d3fcf2b4c77e0cd1c9e2ec3efdc\n" +
	"                    required:\n" +
	"                      - nonce_id\n" +
	"                      - nonce\n" +
	"                      - salt\n" +
	"  /v1/login/input/execute/:\n" +
	"    post:\n" +
	"      summary: Execute a secure login\n" +
	"      tags:\n" +
	"        - Login\n" +
	"      security: []\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                username:\n" +
	"                  $ref: \"#/components/schemas/Username\"\n" +
	"                password:\n" +
	"                  type: string\n" +
	"                  description: \"`Hex(AesGcm256(password))` (with `salt` and `nonce` in a prepare response)\"\n" +
	"                  example: b41a4b0294c07d7ab5b60682a93870cdb0927cf52a76\n" +
	"                nonce_id:\n" +
	"                  $ref: \"#/components/schemas/Code\"\n" +
	"              required:\n" +
	"                - username\n" +
	"                - password\n" +
	"                - nonce_id\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      commit_token:\n" +
	"                        $ref: \"#/components/schemas/JwtToken\"\n" +
	"                    required:\n" +
	"                      - commit_token\n" +
	"  /v1/login/input/commit/:\n" +
	"    post:\n" +
	"      summary: Commit a secure login\n" +
	"      tags:\n" +
	"        - Login\n" +
	"      security:\n" +
	"        - CommitAuth: []\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                device_uid:\n" +
	"                  $ref: \"#/components/schemas/DeviceUID\"\n" +
	"                auth_code:\n" +
	"                  $ref: \"#/components/schemas/AuthCode\"\n" +
	"              required:\n" +
	"                - device_uid\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      access_token:\n" +
	"                        $ref: \"#/components/schemas/JwtToken\"\n" +
	"                      refresh_token:\n" +
	"                        $ref: \"#/components/schemas/JwtToken\"\n" +
	"                    required:\n" +
	"                      - access_token\n" +
	"                      - refresh_token\n" +
	"  /v1/login/refresh/:\n" +
	"    post:\n" +
	"      summary: Refresh tokens\n" +
	"      tags:\n" +
	"        - Login\n" +
	"      security:\n" +
	"        - RefreshAuth: []\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                device_uid:\n" +
	"                  $ref: \"#/components/schemas/DeviceUID\"\n" +
	"                rotate:\n" +
	"                  type: boolean\n" +
	"                  description: TRUE when we need a new Refresh token\n" +
	"              required:\n" +
	"                - device_uid\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      access_token:\n" +
	"                        $ref: \"#/components/schemas/JwtToken\"\n" +
	"                      refresh_token:\n" +
	"                        $ref: \"#/components/schemas/JwtToken\"\n" +
	"                    required:\n" +
	"                      - access_token\n" +
	"  /v1/login/logout/:\n" +
	"    post:\n" +
	"      summary: Logout\n" +
	"      tags:\n" +
	"        - Login\n" +
	"      security:\n" +
	"        - RefreshAuth: []\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  message:\n" +
	"                    $ref: \"#/components/schemas/ResponseMessage\"\n" +
	"  /v1/kyc/init/:\n" +
	"    post:\n" +
	"      summary: Init Kyc\n" +
	"      tags:\n" +
	"        - KYC\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                full_name:\n" +
	"                  type: string\n" +
	"                dob:\n" +
	"                  $ref: \"#/components/schemas/Date\"\n" +
	"                nationality:\n" +
	"                  $ref: \"#/components/schemas/Nationality\"\n" +
	"                residential_address:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"                postal_code:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                city:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"                country:\n" +
	"                  $ref: \"#/components/schemas/Nationality\"\n" +
	"              required:\n" +
	"                - full_name\n" +
	"                - dob\n" +
	"                - nationality\n" +
	"                - residential_address\n" +
	"                - postal_code\n" +
	"                - city\n" +
	"                - country\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      kyc_code:\n" +
	"                        $ref: \"#/components/schemas/Code\"\n" +
	"                      user_code:\n" +
	"                        $ref: \"#/components/schemas/Code\"\n" +
	"                      jumio_api_token:\n" +
	"                        $ref: \"#/components/schemas/Note\"\n" +
	"                      jumio_api_secret:\n" +
	"                        $ref: \"#/components/schemas/Note\"\n" +
	"  /v1/kyc/init/url/:\n" +
	"    post:\n" +
	"      summary: Init Kyc for web\n" +
	"      tags:\n" +
	"        - KYC\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                full_name:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"                dob:\n" +
	"                  $ref: \"#/components/schemas/Date\"\n" +
	"                nationality:\n" +
	"                  $ref: \"#/components/schemas/Nationality\"\n" +
	"                residential_address:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"                postal_code:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                city:\n" +
	"                  $ref: \"#/components/schemas/Note\"\n" +
	"                country:\n" +
	"                  $ref: \"#/components/schemas/Nationality\"\n" +
	"              required:\n" +
	"                - full_name\n" +
	"                - dob\n" +
	"                - nationality\n" +
	"                - residential_address\n" +
	"                - postal_code\n" +
	"                - city\n" +
	"                - country\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      redirect_url:\n" +
	"                        $ref: \"#/components/schemas/Note\"\n" +
	"  /v1/kyc/submit/:\n" +
	"    post:\n" +
	"      summary: Submit kyc\n" +
	"      tags:\n" +
	"        - KYC\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                kyc_code:\n" +
	"                  $ref: \"#/components/schemas/Code\"\n" +
	"                reference:\n" +
	"                  $ref: \"#/components/schemas/Code\"\n" +
	"              required:\n" +
	"                - kyc_code\n" +
	"                - reference\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: string\n" +
	"                    example: null\n" +
	"  /v1/kyc/get/:\n" +
	"    post:\n" +
	"      summary: Get kyc\n" +
	"      tags:\n" +
	"        - KYC\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      request:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          id:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          uid:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          status:\n" +
	"                            $ref: \"#/components/schemas/Status\"\n" +
	"                          full_name:\n" +
	"                            type: string\n" +
	"                          username:\n" +
	"                            $ref: \"#/components/schemas/Username\"\n" +
	"                          note:\n" +
	"                            $ref: \"#/components/schemas/Note\"\n" +
	"                          dob:\n" +
	"                            $ref: \"#/components/schemas/Date\"\n" +
	"                          nationality:\n" +
	"                            $ref: \"#/components/schemas/Nationality\"\n" +
	"                          residential_address:\n" +
	"                            type: string\n" +
	"                          postal_code:\n" +
	"                            type: string\n" +
	"                          city:\n" +
	"                            type: string\n" +
	"                          country:\n" +
	"                            $ref: \"#/components/schemas/Nationality\"\n" +
	"                      user_type:\n" +
	"                        $ref: \"#/components/schemas/UserType\"\n" +
	"                      remaining_attempts:\n" +
	"                        $ref: \"#/components/schemas/Number\"\n" +
	"                      blacklist_info:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          in_blacklist:\n" +
	"                            type: boolean\n" +
	"                          reason:\n" +
	"                            $ref: \"#/components/schemas/Note\"\n" +
	"  /v1/kyc/meta/:\n" +
	"    post:\n" +
	"      summary: Kyc meta\n" +
	"      tags:\n" +
	"        - KYC\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      reminder:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          deadline_time:\n" +
	"                            $ref: \"#/components/schemas/Timestamp\"\n" +
	"                          is_enabled:\n" +
	"                            type: boolean\n" +
	"  /v1/kyc/validate/email/:\n" +
	"    post:\n" +
	"      summary: Kyc validate email\n" +
	"      tags:\n" +
	"        - KYC\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                email:\n" +
	"                  $ref: \"#/components/schemas/Email\"\n" +
	"              required:\n" +
	"                - email\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      is_valid_email:\n" +
	"                        type: boolean\n" +
	"                      message:\n" +
	"                        $ref: \"#/components/schemas/Note\"\n" +
	"  /v1/s2s/kyc/get/:\n" +
	"    post:\n" +
	"      summary: Get kyc\n" +
	"      tags:\n" +
	"        - S2S\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                uid:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"              required:\n" +
	"                - uid\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      request:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          id:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          uid:\n" +
	"                            $ref: \"#/components/schemas/ID\"\n" +
	"                          status:\n" +
	"                            $ref: \"#/components/schemas/Status\"\n" +
	"                          full_name:\n" +
	"                            type: string\n" +
	"                          username:\n" +
	"                            $ref: \"#/components/schemas/Username\"\n" +
	"                          note:\n" +
	"                            $ref: \"#/components/schemas/Note\"\n" +
	"                          dob:\n" +
	"                            $ref: \"#/components/schemas/Date\"\n" +
	"                          nationality:\n" +
	"                            $ref: \"#/components/schemas/Nationality\"\n" +
	"                          residential_address:\n" +
	"                            type: string\n" +
	"                          postal_code:\n" +
	"                            type: string\n" +
	"                          city:\n" +
	"                            type: string\n" +
	"                          country:\n" +
	"                            $ref: \"#/components/schemas/Nationality\"\n" +
	"                      user_type:\n" +
	"                        $ref: \"#/components/schemas/UserType\"\n" +
	"                      remaining_attempts:\n" +
	"                        $ref: \"#/components/schemas/Number\"\n" +
	"                      blacklist_info:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          in_blacklist:\n" +
	"                            type: boolean\n" +
	"                          reason:\n" +
	"                            $ref: \"#/components/schemas/Note\"\n" +
	"  /v1/s2s/kyc/meta/:\n" +
	"    post:\n" +
	"      summary: Kyc meta\n" +
	"      tags:\n" +
	"        - S2S\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      reminder:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          deadline_time:\n" +
	"                            $ref: \"#/components/schemas/Timestamp\"\n" +
	"                          is_enabled:\n" +
	"                            type: boolean\n" +
	"  /v1/s2s/kyc/send/email/:\n" +
	"    post:\n" +
	"      summary: Kyc send email\n" +
	"      tags:\n" +
	"        - S2S\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                request_id:\n" +
	"                  $ref: \"#/components/schemas/ID\"\n" +
	"                request_status:\n" +
	"                  $ref: \"#/components/schemas/Status\"\n" +
	"              required:\n" +
	"                - request_id\n" +
	"                - request_status\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: string\n" +
	"                    example: null\n" +
	"  /v1/s2s/kyc/validate/email/:\n" +
	"    post:\n" +
	"      summary: Kyc validate email\n" +
	"      tags:\n" +
	"        - S2S\n" +
	"      requestBody:\n" +
	"        required: true\n" +
	"        content:\n" +
	"          application/json:\n" +
	"            schema:\n" +
	"              type: object\n" +
	"              properties:\n" +
	"                email:\n" +
	"                  $ref: \"#/components/schemas/Email\"\n" +
	"              required:\n" +
	"                - email\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      is_valid_email:\n" +
	"                        type: boolean\n" +
	"                      message:\n" +
	"                        $ref: \"#/components/schemas/Note\"\n" +
	"  /v1/meta/status/:\n" +
	"    post:\n" +
	"      summary: Meta status\n" +
	"      tags:\n" +
	"        - Meta\n" +
	"      responses:\n" +
	"        \"200\":\n" +
	"          description: OK\n" +
	"          content:\n" +
	"            application/json:\n" +
	"              schema:\n" +
	"                type: object\n" +
	"                properties:\n" +
	"                  errors:\n" +
	"                    $ref: \"#/components/schemas/Errors\"\n" +
	"                  data:\n" +
	"                    type: object\n" +
	"                    properties:\n" +
	"                      registration:\n" +
	"                        type: object\n" +
	"                        properties:\n" +
	"                          banned_info:\n" +
	"                            type: object\n" +
	"                            properties:\n" +
	"                              ip:\n" +
	"                                $ref: \"#/components/schemas/IP\"\n" +
	"                              country:\n" +
	"                                $ref: \"#/components/schemas/Country\"\n" +
	"                              message:\n" +
	"                                $ref: \"#/components/schemas/Note\"\n" +
	"                          is_banned:\n" +
	"                            type: boolean\n" +
	""
