package kycmod

var kycContentHaveFund = `
	<tr>
		<td class="content">
			To prevent unnecessary delays with your KYC verification,
			please withdraw any funds stored in your Torque Accounts you have created after 23:59:59hrs GMT +7 on 5 September 2020 (“Restricted Accounts”) within seven (7) days from the date of this email.
		</td>
	</tr>
	<tr><td>&nbsp;</td></tr>
	<tr>
		<td class="content">
			If, through your inaction, we are still unable to process your KYC verification after the 7 days, we may do one or more of the following:
			(a) suspend your access to our services for all your Torque Accounts and limiting your account functions to "Withdrawal" and "View Balance"; and
			(b) prevent you from completing any actions via Torque, including funding any of your Torque accounts and/or receiving any Rewards.
		<td>
	</tr>
	<tr><td>&nbsp;</td></tr>
`

var kycContentNoFund = `
	<tr><td>&nbsp;</td></tr>
	<tr>
		<td class="content">
			To prevent unnecessary delays with your KYC verification,
			we will proceed to close any Torque  Accounts you have created after 23:59:59hrs GMT +7 on 5 September 2020 (“Restricted Accounts”) which do not meet our Minimum Balance requirements.
		<td>
	</tr>
	<tr><td>&nbsp;</td></tr>
`

var kycContentHaveAndNoFund = `
	<tr>
		<td class="content">
			To prevent unnecessary delays with your KYC verification,
			you are required to withdraw any funds stored in your Torque Accounts you have created after 23:59:59hrs GMT +7 on 5 September 2020 (“Restricted Accounts”) within 7 days from the date of this email.
		</td>
	</tr>
	<tr><td>&nbsp;</td></tr>
	<tr>
		<td class="content">
			If, through your inaction, we are still unable to process your KYC verification after the 7 days, we may do one or more of the following:
		</td>
	</tr>
	<tr>
		<td class="content">
			(a) suspend your access to our services for all your Torque Accounts and limit your account functions to "Withdrawal" and "View Balance";
		</td>
	</tr>
	<tr>
		<td class="content">
			(b) prevent you from completing any actions via Torque, including funding any of your Torque accounts and/or receiving any Rewards.
		</td>
	</tr>
	<tr><td>&nbsp;</td></tr>
	<tr>
		<td class="content">
			Accordingly, we will proceed to close your Restricted Torque Accounts, including the accounts with any residual funds which are less or equal to the network fees.
		</td>
	</tr>
	<tr><td>&nbsp;</td></tr>
`
