package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type Player struct {
	Name      string
	Nat20     int
	Nat1      int
	DanoTotal int
	DanoMax   int
	CuraTotal int
	CuraMax   int
	Quedas    int
	Mortes    int
	Custom    string
}

var (
	players = []string{"Gustavo", "Mariana", "Pedro", "Joao", "Janis", "Catti", "Maria", "Eric", "Andre"}
	data    = make(map[string]*Player)
	mu      sync.Mutex
	active  = ""
)

func init() {
	for _, name := range players {
		data[name] = &Player{Name: name}
	}
}

func renderTable() string {
	var sb strings.Builder
	sb.WriteString("```ansi\n")
	sb.WriteString("\x1b[1;34m==================================================\x1b[0m\n")
	sb.WriteString("           \U0001f4d6 \x1b[1;37mGRIMOIRE: AUTOS DA AVENTURA\x1b[0m\n")
	sb.WriteString("\x1b[1;34m==================================================\x1b[0m\n")
	sb.WriteString("\x1b[1;33mJOGADOR  | N20 | N1 | D.TOTAL | D.MAX | C.TOTAL | C.MAX | Q | M \x1b[0m\n")
	sb.WriteString("--------------------------------------------------\n")

	for _, name := range players {
		p := data[name]
		row := fmt.Sprintf("%-8s | %-3d | %-2d | %-7d | %-5d | %-7d | %-5d | %-1d | %-1d\n",
			p.Name, p.Nat20, p.Nat1, p.DanoTotal, p.DanoMax, p.CuraTotal, p.CuraMax, p.Quedas, p.Mortes)
		if p.Custom != "" {
			row += fmt.Sprintf(" \u2514\u2500 \x1b[0;32m%s\x1b[0m\n", p.Custom)
		}
		sb.WriteString(row)
	}
	sb.WriteString("--------------------------------------------------\n")
	if active != "" {
		sb.WriteString(fmt.Sprintf("\x1b[1;32mFoco Atual: %s\x1b[0m\n", active))
	}
	sb.WriteString("```")
	return sb.String()
}

func main() {
	_ = godotenv.Load()
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN is not set")
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Erro ao criar sess\u00e3o:", err)
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if i.ApplicationCommandData().Name == "grimoire" {
				respondDashboard(s, i)
			}

		case discordgo.InteractionMessageComponent:
			handleComponents(s, i)

		case discordgo.InteractionModalSubmit:
			handleModals(s, i)
		}
	})

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{Name: "grimoire", Description: "Abre o painel Grimoire"})
		fmt.Println("Grimoire Online! Use /grimoire no Discord.")
	})

	err = dg.Open()
	if err != nil {
		log.Fatal("Erro ao abrir conex\u00e3o:", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}

func respondDashboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    renderTable(),
			Components: createComponents(),
		},
	})
}

func createComponents() []discordgo.MessageComponent {
	var options []discordgo.SelectMenuOption
	for _, name := range players {
		options = append(options, discordgo.SelectMenuOption{Label: name, Value: name})
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{CustomID: "select_player", Placeholder: "\U0001f464 Selecionar Jogador", Options: options},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "N20", CustomID: "add_n20", Style: discordgo.SuccessButton},
			discordgo.Button{Label: "N1", CustomID: "add_n1", Style: discordgo.DangerButton},
			discordgo.Button{Label: "Queda", CustomID: "add_q", Style: discordgo.SecondaryButton},
			discordgo.Button{Label: "Morte", CustomID: "add_m", Style: discordgo.SecondaryButton},
		}},
		discordgo.ActionsRow{Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "\U0001f4dd Registrar Dano/Cura", CustomID: "open_modal", Style: discordgo.PrimaryButton},
			discordgo.Button{Label: "\u2699\ufe0f Custom", CustomID: "open_custom", Style: discordgo.SecondaryButton},
		}},
	}
}

func handleComponents(s *discordgo.Session, i *discordgo.InteractionCreate) {
	mu.Lock()
	defer mu.Unlock()

	id := i.MessageComponentData().CustomID

	if id == "select_player" {
		active = i.MessageComponentData().Values[0]
	} else if active == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "Selecione um jogador primeiro!", Flags: 64},
		})
		return
	}

	p := data[active]
	switch id {
	case "add_n20":
		p.Nat20++
	case "add_n1":
		p.Nat1++
	case "add_q":
		p.Quedas++
	case "add_m":
		p.Mortes++
	case "open_modal":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "modal_data", Title: "Registrar para " + active,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{CustomID: "val_dano_total", Label: "Valor de Dano Total", Style: discordgo.TextInputShort, Placeholder: "0", Value: strconv.Itoa(p.DanoTotal)},
					}},
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{CustomID: "val_dano_max", Label: "Valor de Dano Maximo", Style: discordgo.TextInputShort, Placeholder: "0", Value: strconv.Itoa(p.DanoMax)},
					}},
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{CustomID: "val_cura_total", Label: "Valor de Cura Total", Style: discordgo.TextInputShort, Placeholder: "0", Value: strconv.Itoa(p.CuraTotal)},
					}},
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{CustomID: "val_cura_max", Label: "Valor de Cura Maximo", Style: discordgo.TextInputShort, Placeholder: "0", Value: strconv.Itoa(p.CuraMax)},
					}},
				},
			},
		})
		return
	case "open_custom":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: "modal_custom", Title: "Anota\u00e7\u00e3o Custom: " + active,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{CustomID: "val_custom", Label: "Texto (ex: Sorte: 2)", Style: discordgo.TextInputShort, Value: p.Custom},
					}},
				},
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{Content: renderTable(), Components: createComponents()},
	})
}

func handleModals(s *discordgo.Session, i *discordgo.InteractionCreate) {
	mu.Lock()
	defer mu.Unlock()
	p := data[active]
	d := i.ModalSubmitData()

	if d.CustomID == "modal_data" {
		dano_total, _ := strconv.Atoi(d.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
		dano_max, _ := strconv.Atoi(d.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
		cura_total, _ := strconv.Atoi(d.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
		cura_max, _ := strconv.Atoi(d.Components[3].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)

		p.DanoTotal = dano_total
		p.DanoMax = dano_max
		p.CuraTotal = cura_total
		p.CuraMax = cura_max
	} else {
		p.Custom = d.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{Content: renderTable(), Components: createComponents()},
	})
}
