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

	mg := mailgun.NewMailgun("mail.kingscollegecam.com", os.Getenv("MailGunKey"), os.Getenv("MailGunPubKey"))

	message := mg.NewMessage(
		/* Sender */ "noreply@store.kings.cam.ac.uk",
		/* Subject*/ "King's College Chapel Ticket",
		/* Body */ "",
		/* Recipient */ booking.Email)
	message.SetHtml("<html><p>Booking reference: " + booking.UUID + "<br/>Booking date and sessions: " + booking.Date + " " + booking.Session + "<br/>Total for " + strconv.Itoa(booking.Ntickets) + " tickets (" + strconv.Itoa(booking.NAdults) + " adults, " + strconv.Itoa(booking.Nconcession) + " students and " + strconv.Itoa(booking.Nchild) + " child) and " + strconv.Itoa(booking.Nguides) + " guides is: <span>&#163;</span>" + fmt.Sprintf("%.2f", booking.Total) + "</p><p>Please present your ticket at the main gate on King’s Parade and you will be directed to the entrance.<img src=\"cid:" + booking.UUID + ".png\"></p><p>Find out how to get here at <a href=\"http://www.kings.cam.ac.uk/visit/getting-to-kings.html\">http://www.kings.cam.ac.uk/visit/getting-to-kings.html</a></p><p>Whilst at King's, do visit the King's College Visitor Centre on King's Parade (opposite the main gate) for souvenirs, gifts, CDs and clothing. Alternatively visit the online shop at <a href=\"http://shop.kings.cam.ac.uk\">http://shop.kings.cam.ac.uk</a></p></html>")

	message.AddInline(booking.UUID + ".png")
	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
