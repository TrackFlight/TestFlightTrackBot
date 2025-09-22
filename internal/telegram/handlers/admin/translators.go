package admin

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/GoBotApiOfficial/gobotapi/parser"
	"github.com/GoBotApiOfficial/gobotapi/types"
	"github.com/TrackFlight/TestFlightTrackBot/internal/telegram/core"
	"github.com/TrackFlight/TestFlightTrackBot/internal/translator"
)

func EditVarAction(ctx *core.UpdateContext, message types.Message) error {
	parsedText := parser.Parse(
		message.Text,
		message.Entities,
	)
	parsedText = strings.ReplaceAll(parsedText, "&lt;", "<")
	parsedText = strings.ReplaceAll(parsedText, "&gt;", ">")
	dataWithHtml := strings.SplitN(parsedText, " ", 3)
	data := strings.SplitN(message.Text, " ", 3)
	var errUsage, failUpdate, successUpdate translator.Key
	var action translator.Action
	requiredParams := 3
	switch data[0] {
	case "/new_var":
		errUsage = translator.NewVarErrorUsage
		failUpdate = translator.NewVarAlreadyExists
		successUpdate = translator.NewVarSuccess
		action = translator.Add
	case "/del_var":
		errUsage = translator.RemoveVarErrorUsage
		failUpdate = translator.RemoveOrEditVarNotFound
		successUpdate = translator.RemoveVarSuccess
		action = translator.Remove
		requiredParams = 2
	case "/edit_var":
		errUsage = translator.EditVarErrorUsage
		failUpdate = translator.RemoveOrEditVarNotFound
		successUpdate = translator.EditVarSuccess
		action = translator.Update
	}
	if len(data) < requiredParams {
		return ctx.SendMessage(
			message.Chat.ID,
			ctx.Translator.T(errUsage),
		)
	}
	varName := strings.ToUpper(data[1])
	var value string
	if action == translator.Update || action == translator.Add {
		value = dataWithHtml[2]
	}
	if ok, err := translator.Edit(varName, value, action); err != nil {
		return err
	} else if !ok {
		return ctx.SendMessage(
			message.Chat.ID,
			ctx.Translator.TWithData(
				failUpdate,
				map[string]string{
					"VarName": varName,
				},
			),
		)
	} else {
		return ctx.SendMessage(
			message.Chat.ID,
			ctx.Translator.T(successUpdate),
		)
	}
}

func SearchVar(ctx *core.UpdateContext, message types.Message) error {
	data := strings.SplitN(message.Text, " ", 2)
	if len(data) < 2 {
		return ctx.SendMessage(
			message.Chat.ID,
			ctx.Translator.T(translator.SearchVarErrorUsage),
		)
	}
	varName := strings.ToUpper(data[1])
	if v, err := translator.SearchVar(varName); err != nil {
		return err
	} else if len(v) == 0 {
		return ctx.SendMessage(
			message.Chat.ID,
			ctx.Translator.TWithData(
				translator.SearchVarNotFound,
				map[string]string{
					"Text": varName,
				},
			),
		)
	} else {
		var builder strings.Builder
		for key, value := range v {
			escapedValue := html.EscapeString(value)
			segment := fmt.Sprintf(
				"<b>%s</b>: \n<code>%s</code>\n\n",
				key,
				escapedValue,
			)

			if builder.Len()+len(segment) >= 4096 {
				break
			}
			builder.WriteString(segment)
		}
		return ctx.SendMessage(
			message.Chat.ID,
			ctx.Translator.TWithData(
				translator.SearchVarSuccess,
				map[string]string{
					"Results":      builder.String(),
					"ResultsCount": strconv.Itoa(len(v)),
				},
			),
		)
	}
}
