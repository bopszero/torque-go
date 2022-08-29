package kycmod

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comcontext"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/auth/emailmod"
	"gitlab.com/snap-clickstaff/torque-go/lib/utils"
)

func ParseEmailOriginal(ctx comcontext.Context, email string) string {
	email = strings.ToLower(email)
	emailParts := strings.Split(email, "@")
	if len(emailParts) != 2 {
		return email
	}
	var (
		emailName   = emailParts[0]
		emailDomain = emailParts[1]
	)
	switch {
	case isMatchEmailDomainPattern(ctx, EmailDomainYahooPattern, emailDomain):
		dashSignIndex := strings.Index(emailName, "-")
		if dashSignIndex != -1 {
			emailName = emailName[:dashSignIndex]
		}
		break
	case isMatchEmailDomainPattern(ctx, EmailDomainGmailPattern, emailDomain),
		isMatchEmailDomainPattern(ctx, EmailDomainProtonMailPattern, emailDomain):
		plusSignIndex := strings.Index(emailName, "+")
		if plusSignIndex != -1 {
			emailName = emailName[:plusSignIndex]
		}
		emailName = strings.ReplaceAll(emailName, ".", "")
		break
	default:
		break
	}
	return fmt.Sprintf("%s@%s", emailName, emailDomain)
}

func isMatchEmailDomainPattern(ctx comcontext.Context, regexPattern string, userEmailDomain string) bool {
	matched, err := regexp.Match(regexPattern, []byte(userEmailDomain))
	if err != nil {
		comlogging.GetLogger().
			WithType(constants.LogTypeKYC).
			WithContext(ctx).
			WithError(err).
			WithFields(logrus.Fields{
				"regex_pattern":     regexPattern,
				"user_email_domain": userEmailDomain,
			}).
			Errorf("regex pattern has error | err=%v", err.Error())
	}
	return matched
}

func genKycRejectVerificationStatus(verificationStatus string, rejectCode interface{}) string {
	rejectCodeStr := comutils.Stringify(rejectCode)
	return fmt.Sprintf(
		"%v:%v:%v",
		KeyVerificationStatusJumio,
		strings.ToLower(verificationStatus),
		strings.ToLower(rejectCodeStr),
	)
}

func TruncateName(name string) string {
	name = FullNameDirtyCharacterRegex.ReplaceAllString(name, "")
	var (
		nameParts  = strings.Split(name, " ")
		validParts = make([]string, 0, len(nameParts))
	)
	for _, part := range nameParts {
		if part != "" {
			validParts = append(validParts, part)
		}
	}

	return strings.Join(validParts, "")
}

func isPreferLeftRequest(leftReq, rightReq models.KycRequest) bool {
	if leftReq.Status == rightReq.Status {
		return leftReq.CreateTime > rightReq.CreateTime
	}
	var (
		leftPriority  = RequestStatusPriorityMap[leftReq.Status]
		rightPriority = RequestStatusPriorityMap[rightReq.Status]
	)
	if leftPriority == 0 {
		return false
	}
	if rightPriority == 0 {
		return true
	}
	return leftPriority < rightPriority
}

func renderEmailTemplate(templatePath string, templateData interface{}) (_ string, err error) {
	emailTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		err = utils.WrapError(err)
		return
	}

	var buf bytes.Buffer
	if err = emailTemplate.Execute(&buf, templateData); err != nil {
		err = utils.WrapError(err)
		return
	}

	return buf.String(), nil
}

func sendEmailTemplate(
	ctx comcontext.Context,
	emailAddress, subject string, templatePath string, templateData interface{},
) error {
	body, err := renderEmailTemplate(templatePath, templateData)
	if err != nil {
		return err
	}

	email := emailmod.NewEmail().
		SetSubject(subject).
		AddToAddress(emailAddress)
	email.SetBody(echo.MIMETextHTMLCharsetUTF8, body)

	return emailmod.Send(ctx, email)
}

func isMatchDocumentNationality(ctx comcontext.Context, kycRequestNationality string, document Document) bool {
	if document.Nationality == "" && document.IssuingCountry == "" {
		return false
	}
	var documentNationality string
	if document.Nationality != "" {
		documentNationality = document.Nationality
	} else {
		documentNationality = document.IssuingCountry
	}
	var country models.Country
	err := database.GetDbSlave().
		First(&country, &models.Country{CodeIso2: kycRequestNationality, CodeIso3: documentNationality}).
		Error
	if err != nil {
		comlogging.GetLogger().
			WithType(constants.LogTypeKYC).
			WithContext(ctx).
			WithError(err).
			WithFields(logrus.Fields{
				"document_nationality": documentNationality,
				"request_nationality":  kycRequestNationality,
			}).
			Error("find document nationality failed | err=%s", err.Error())
		return false
	}
	return true
}
