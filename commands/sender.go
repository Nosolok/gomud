package commands

import (
	"github.com/Vehsamrak/gomud/console"
	"net"
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
)

const DEFAULT_ENCODING = "utf-8"
const ENCODING_WINDOWS_1251 = "windows-1251"
const ENCODING_KOI8 = "koi8-r"
const ENCODING_WINDOWS_1252 = "windows-1252"

type Sender struct {
	ConnectionPointer *net.Conn
	charset string
}

func (sender *Sender) toClient(message string)  {
	console.Client(sender.ConnectionPointer, sender.encodeToClientCharset(message))
}

func (sender *Sender) toServer(message string)  {
	console.Server(message)
}

func (sender *Sender) encodeToClientCharset(message string) string {
	processedMessageBytes := sender.fixYaLetter(message)
	charsetTranslator, _ := charset.TranslatorTo(sender.charset)
	_, translatedMessageBytes, _ := charsetTranslator.Translate(processedMessageBytes, false)

	return string(translatedMessageBytes)
}

func (sender *Sender) encodeToUtf8(message string) string {
	charsetTranslator, _ := charset.TranslatorFrom(sender.charset)
	_, translatedMessageBytes, _ := charsetTranslator.Translate([]byte(message), false)

	return string(translatedMessageBytes)
}

func (sender *Sender) fixYaLetter(message string) []byte {
	messageBytes := []byte(message)

	if sender.charset != ENCODING_WINDOWS_1251 {
		return messageBytes
	}

	var processedMessageBytes []byte
	prebyte := false

	// double "я"-letter (209 & 143 bytes) to fix CP1251 issue
	for _, messageByte := range messageBytes {
		if prebyte && messageByte == 143 {
			processedMessageBytes = append(processedMessageBytes, 143)
			processedMessageBytes = append(processedMessageBytes, 209)
		}

		if messageByte == 209 {
			prebyte = true
		} else {
			prebyte = false
		}

		processedMessageBytes = append(processedMessageBytes, messageByte)
	}

	return processedMessageBytes
}