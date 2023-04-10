package infinity

type EmailService interface {
	SendEmailFromTemplateFile(string, string, string)
	SendAccountCreationEmail(*User)
	SendPasswordResetEmail(*User, string)
	SendPaymentPlanExpiryReminderEmail(*User, uint)
	SendPaymentPlanExpiredEmail(*User)
	SendAccountCancellationEmail(*User)
	SendPurchasedRegularPremiumEmail(*User, *PaymentPlan, *PaymentTransaction)
	SendPurchasedTierTwoPremiumEmail(*User)
}
