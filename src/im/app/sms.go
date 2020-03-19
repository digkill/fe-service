package app

import (
	"fmt"
	"im/mlog"
	"im/model"
	smsaero "im/utils"
)

func (a *App) SendVerifySms(phone, locale, msg string) *model.AppError {

	message := model.SmsNotification{
		Phone:   phone,
		Message: msg,
	}

	a.Srv.Go(func() {
		print(message.Message)
		fmt.Println(message)
		//a.sendToSmsProxy(message)
		if !*a.Config().ServiceSettings.EnableDeveloper {
			if app, err := a.GetApplication(locale); err != nil {
				mlog.Error(err.Error())
			} else {
				a.sendToSmsProxy(app.SmsLogin, app.SmsApiKey, message)
			}
		}
	})

	return nil
}

func (a *App) sendToSmsProxy(login, apiKey string, msg model.SmsNotification) {

	//client := smsaero.NewClient("ivan@russianit.ru", "mmRjD5mOoMkvVsuIAMiVwX6i9czQ", "", "")
	//client := smsaero.NewClient("amgraal@mail.ru", "QxXk6d6055rG7bgttUUeQEtHrsmn", "", "")
	client := smsaero.NewClient(login, apiKey, "", "")
	msgAero := smsaero.MessageRequest{
		Numbers: []string{msg.Phone},
		Sign:    "SMS Aero",
		Text:    msg.Message,
		Channel: smsaero.ChannelDirect,
	}
	resp, err := client.Send(msgAero)
	if err != nil {
		mlog.Error(err.Error())
	}

	fmt.Printf("%#v\n", resp)

	return
}
