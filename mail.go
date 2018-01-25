package tickets

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/mailgun/mailgun-go.v1"
)

// Sends email
func sendmail(booking *Booking) {

	mg := mailgun.NewMailgun("sandbox58b570d32fb64c16a83fd3f306d57045.mailgun.org", os.Getenv("MailGunKey"), os.Getenv("MailGunPubKey"))

	message := mg.NewMessage(
		/* Sender */ "store@kings.cam.ac.uk",
		/* Subject*/ "King's College Chapel Ticket",
		/* Body */ "",
		/* Recipient */ booking.Email)
	message.SetHtml("<html><p>Booking reference: " + booking.UUID + "<br/>Booking date and sessions: " + booking.Date + " " + booking.Session + "<br/>Total for " + strconv.Itoa(booking.Ntickets) + " tickets (" + strconv.Itoa(booking.NAdults) + " adults, " + strconv.Itoa(booking.Nconcession) + " students and " + strconv.Itoa(booking.Nchild) + " child) and " + strconv.Itoa(booking.Nguides) + "guides is: <span>&#163;</span>" + fmt.Sprintf("%.2f", booking.Total) + "</p><p>Whilst at King's, do visit the King's College Visitor Centre on King's Parade (opposite the main gate) for souvenirs, gifts, CDs and clothing. On buy online  at <a href=\"http://shop.kings.cam.ac.uk\">http://shop.kings.cam.ac.uk</a></p></html>")

	message.AddAttachment(booking.UUID + ".png")
	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
