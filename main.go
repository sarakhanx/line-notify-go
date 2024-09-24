package main

import (
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

// LineNotifyURL is the API endpoint for Line Notify
const LineNotifyURL = "https://notify-api.line.me/api/notify"

// Set client request to 3rd party
func sendLineNotification(token, message string) error {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+token).
		SetFormData(map[string]string{
			// "message" ก้อนนี้ สามารถเปลี่ยนเป็นอะไรก็ได้ เพื่อเอาไว้ตอบกลับไปทาง Notify
			"message": message,
		}).
		Post(LineNotifyURL)

	if err != nil {
		return err
	}

	if resp.StatusCode() != fiber.StatusOK {
		return fmt.Errorf("failed to send notification: %s", resp.Status())
	}

	log.Printf("Notification sent successfully. Response: %s", resp)
	return nil
}

func lineNotifyHandler(c *fiber.Ctx) error {
	token := "ArGLXDLMy5kIi3pLCT0Yhnbd8FghZ0MaSZPB2H54qDG" //you can try if you want but it's will notice me ^^
	type Body struct {
		Message string `json:"message"`
		// StickerPackageId string `json:"stickerPackageId"`
		// StickerId        string `json:"stickerId"`
	}

	body := Body{}
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Error reading body: %v", err))
	}
	message := body.Message

	//EXPLAIN - เรียกใช้ฟังก์ชันสำหรับส่งข้อความไปยัง Line Notify(set client request to 3rd party)
	err = sendLineNotification(token, message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error sending notification: %v", err))
	}

	return c.SendString("Line Notify received")
}

func main() {
	app := fiber.New()

	app.Post("/line-notify", lineNotifyHandler)

	fmt.Println("Starting server on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
