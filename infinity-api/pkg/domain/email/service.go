package email

import (
	"fmt"

	"git.misc.vee.bz/carnagel/infinity-api/pkg"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"os"
)

const footerTemplate = `
<a href="https://camtube.co" target="_blank" style="display: inline-block;">
<svg version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px"
	 width="250px" height="102.87px" viewBox="0 0 552.662 102.87" enable-background="new 0 0 552.662 102.87"
	 xml:space="preserve">
<g>
<g>
       <g>
		<path fill="#FB5255" d="M482.296,66.664c0,11.782,4.433,16.156,10.79,16.156c6.532,0,10.79-4.433,10.79-16.214
			c0-11.782-4.258-16.273-10.79-16.273S482.296,54.883,482.296,66.664z M51.018,69.801c-4.196,3.423-8.834,6.404-15.791,6.404
			c-9.497,0-16.123-6.221-16.123-24.552c0-18.662,6.957-24.883,16.564-24.883c5.301,0,9.718,1.988,13.693,5.521l9.055-12.221
			c-6.294-5.3-13.473-8.503-23.853-8.503C14.908,11.568,0,26.034,0,51.653C0,77.936,13.914,91.96,34.233,91.96
			c11.374,0,19.767-4.638,25.067-9.938L51.018,69.801z M366.664,66.082c0-7.95-4.97-10.711-12.589-10.711h-8.393v20.65h7.619
			C360.369,76.021,366.664,74.365,366.664,66.082z M81.71,13.556L58.189,90.082h19.325l3.976-16.453h22.196l3.975,16.453h19.657
			l-23.411-76.526H81.71z M84.361,60.377l7.951-32.355h0.552l7.841,32.355H84.361z M552.662,102.87l-43.447-42.444
			c0.352,1.89,0.553,3.938,0.553,6.181c0,13.415-6.823,20.764-16.681,20.764c-9.856,0-16.681-7.058-16.681-20.706
			c0-13.414,6.824-20.88,16.681-20.88c0.407,0,0.801,0.031,1.197,0.056L447.361,0h-243v12.746h74.628v50.356
			c0,9.717,3.534,14.466,11.705,14.466c8.172,0,11.706-4.749,11.706-14.466V12.746h18.331v50.908
			c0,15.681-10.049,27.497-30.037,27.497c-20.098,0-30.036-11.154-30.036-27.497V26.77h-20.092v62.502h-18.332V26.77h-17.276
			l-1.011-13.214h-23.08l-10.822,51.46h-0.441l-11.485-51.46h-23.08l-5.742,76.526h17.889l1.436-35.227
			c0.331-8.834,0.11-13.691-0.552-19.875h0.442l12.147,47.262h17.558l11.595-47.262h0.441c-0.331,5.411-0.331,9.275-0.11,17.446
			c0.332,12.589,0.994,25.068,1.656,37.656h12.56v12.788H552.662z M460.661,45.784c5.133,0,8.165,1.4,11.315,3.966l-3.033,3.558
			c-2.45-2.042-5.074-2.975-7.933-2.975c-5.949,0-11.14,4.141-11.14,16.214c0,11.606,4.841,16.156,11.082,16.156
			c4.199,0,6.766-1.691,8.748-3.441l2.858,3.616c-2.508,2.392-6.241,4.491-11.84,4.491c-9.683,0-16.739-7.29-16.739-20.822
			C443.98,53.017,451.678,45.784,460.661,45.784z M440.48,79.263c2.333,0,4.083,1.808,4.083,4.024c0,2.275-1.75,4.083-4.083,4.083
			c-2.274,0-4.024-1.808-4.024-4.083C436.456,81.07,438.206,79.263,440.48,79.263z M389.515,12.746h42.293l-0.257,13.031h-23.705
			v18.221h23.521v12.81h-23.521v19.215h23.962v13.251h-42.293V12.746z M327.351,12.746h25.619c15.901,0,29.484,4.638,29.484,19.546
			c0,9.166-7.067,14.687-13.251,16.123v0.441c7.509,1.546,16.453,5.963,16.453,18c0,17.669-14.907,22.417-31.582,22.417h-26.724
			V12.746z M363.792,34.168c0-5.963-3.645-8.502-11.595-8.502h-6.516v17.668h7.067C360.48,43.334,363.792,40.574,363.792,34.168z"/>
	</g>
</g>
</div>
</svg>
</a>
`

const accountActivationTemplate = `
<p>Welcome %s,</p>
<p>Your account has now been created.</p>
<p>This email is automatically generated and cannot be answered.</p>
`

const paymentPlanExpiredTemplate = `
<p>Hi %s,</p>
<p>
Your current payment plan has expired, press <a target="_blank" href="https://camtube.co/payment-plans">here</a>
to upgrade your payment plan again.
</p>
<p>This email is automatically generated and cannot be answered.</p>
`

const purchasedRegularPremiumTemplate = `
<p>Hi %s,</p>
<p>Your payment plan has been upgraded, details:</p>
<p>%.2f USD / %f Bitcoin</p>
<p>%d days added to your subscription.</p>
<p>This email is automatically generated and cannot be answered.</p>
`

const paymentPlanExpiryReminderTemplate = `
<p>Hi %s,</p>
<p>Your current payment plan expires in %d %s:</p>
<p>This email is automatically generated and cannot be answered.</p>
`

const passwordResetTemplate = `
<p>Hi %s,</p>
<p>A password reset has been initiated on your account, to continue press
<a target="_blank" href="https://camtube.co/password-reset?token=%s">here</a>.
</p>
<p>This email is automatically generated and cannot be answered.</p>
`

type emailService struct {
	app         *infinity.Application
	dialer      *gomail.Dialer
	fromAddress string
}

func NewEmailService(app *infinity.Application, host string, port int, username string, password string, fromAddress string) infinity.EmailService {
	dialer := gomail.NewPlainDialer(host, port, username, password)

	return &emailService{
		app:         app,
		dialer:      dialer,
		fromAddress: fromAddress,
	}
}

func (service *emailService) SendEmailFromTemplateFile(email string, subject, path string) {
	logger := service.app.Logger.WithFields(logrus.Fields{"method": "SendEmailFromTemplateFile", "email": email})
	logger.Debug("Sending email from template file")

	file, err := os.Open(path)
	if err != nil {
		logger.WithError(err).Error("Failed to open file")

		return
	}

	defer file.Close()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		logger.WithError(err).Error("Failed to read body of file")

		return
	}

	bodyWithFooter := `
	%s
	%s
	`

	// Can't send this from a goroutine since it's used by the email-client, causing child routine to terminate
	// before the email has been sent
	m := gomail.NewMessage()
	m.SetHeader("From", service.fromAddress)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", fmt.Sprintf(bodyWithFooter, string(body), footerTemplate))

	if err := service.dialer.DialAndSend(m); err != nil {
		service.app.Logger.WithField("subject", subject).WithError(err).Error("Failed to send email")
	}
}

func (service *emailService) SendAccountCreationEmail(user *infinity.User) {
	service.sendEmail(
		user.Email,
		"[Camtube] Account created",
		service.insertParamsIntoTemplate(accountActivationTemplate, user.Username),
	)
}

func (service *emailService) SendPasswordResetEmail(user *infinity.User, token string) {
	service.sendEmail(
		user.Email,
		"[Camtube] Password reset",
		service.insertParamsIntoTemplate(passwordResetTemplate, user.Username, token),
	)
}

func (service *emailService) SendPaymentPlanExpiryReminderEmail(user *infinity.User, daysUntilExpiration uint) {
	t := daysUntilExpiration
	unit := "days"
	if daysUntilExpiration == 1 {
		t = 24
		unit = "hours"
	}

	service.sendEmail(
		user.Email,
		"[Camtube] Payment plan expiry notice",
		service.insertParamsIntoTemplate(paymentPlanExpiryReminderTemplate, user.Username, t, unit),
	)
}

func (service *emailService) SendPaymentPlanExpiredEmail(user *infinity.User) {
	service.sendEmail(
		user.Email,
		"[Camtube] Payment plan expired",
		service.insertParamsIntoTemplate(paymentPlanExpiredTemplate, user.Username),
	)
}

func (service *emailService) SendAccountCancellationEmail(user *infinity.User) {
}

func (service *emailService) SendPurchasedRegularPremiumEmail(user *infinity.User, plan *infinity.PaymentPlan, transaction *infinity.PaymentTransaction) {
	service.sendEmail(
		user.Email,
		"[Camtube] Payment plan upgraded",
		service.insertParamsIntoTemplate(
			purchasedRegularPremiumTemplate,
			user.Username,
			plan.Price,
			float64(transaction.GetTotalAmountReceivedInSatoshis())/float64(infinity.SatoshisPerBitCoin),
			plan.Duration,
		),
	)
}

func (service *emailService) SendPurchasedTierTwoPremiumEmail(user *infinity.User) {
}

func (service *emailService) sendEmail(toAddress string, subject string, body string) {
	bodyWithFooter := `
	%s
	%s
	`

	m := gomail.NewMessage()
	m.SetHeader("From", service.fromAddress)
	m.SetHeader("To", toAddress)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", fmt.Sprintf(bodyWithFooter, body, footerTemplate))

	go func() {
		if err := service.dialer.DialAndSend(m); err != nil {
			service.app.Logger.WithField("subject", subject).WithError(err).Error("Failed to send email")
		}
	}()
}

func (service *emailService) insertParamsIntoTemplate(template string, params ...interface{}) string {
	return fmt.Sprintf(template, params...)
}
