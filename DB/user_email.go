package DB

import (
  "time"
  
  "bytes"
  "encoding/hex"
  "html/template"
  "net/smtp"

  "go.mongodb.org/mongo-driver/bson"
)

/// blocking send email function
func internalSendConfirmationEmail(user *CreatableUser) error {
  const subject = "Subject: Confirmation for participation in Syrinx\n"
  const mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

  tmpl, err := template.ParseFiles("email_template.html")
  if err != nil {
    return err
  }

  var body bytes.Buffer
  err = tmpl.Execute(&body, struct {
    JoinedMember string
    Email        string
    TeamName     string
    TeamID       string
  }{
    JoinedMember: user.Username,
    Email:        user.Email,
    TeamName:     *user.TeamName,
    TeamID:       hex.EncodeToString(user.TeamID[:]),
  })
  if err != nil {
    return err
  }

  message := subject + mime + body.String()

  err = smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", EMAIL_SENDER, EMAIL_SENDER_PASSWORD, "smtp.gmail.com"),
                      EMAIL_SENDER, []string{user.Email}, []byte(message),
  )
  if err != nil {
    return err
  }

  return nil
}

func internalUpdateEmailStatus(user *CreatableUser) error {
  _, err := UserDB.Coll.UpdateOne(UserDB.Context, bson.M{"user": user.Username}, bson.M{"$set": bson.M{"mailReceived": true}})
  return err
}

/// Must run this as Async
func sendEmailAsync(user *CreatableUser) {
  const maxCount = 60
  var err error
  count := 0

  err = internalSendConfirmationEmail(user)
  for err != nil && count < maxCount {
    count += 1
    time.Sleep(time.Minute)
    err = internalSendConfirmationEmail(user)
  }

  err = internalUpdateEmailStatus(user)
  for err != nil && count < maxCount {
    count += 1
    time.Sleep(time.Minute)
    err = internalUpdateEmailStatus(user)
  }
}
