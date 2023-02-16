package integrate

import (
	"fmt"

	"github.com/NordGus/gokedex/pkg/extract"
)

type PokemonPage struct {
	Parent     NotionPageParent            `json:"parent"`
	Icon       NotionPageExternalFileIcon  `json:"icon"`
	Cover      NotionPageExternalFileCover `json:"cover"`
	Properties PokemonPageProperties       `json:"properties"`
	Children   []NotionPageBlockObject     `json:"children"`
}

type NotionPageParent struct {
	Type       string `json:"type"`
	DatabaseID string `json:"database_id"`
}

type NotionPageExternalFileIcon struct {
	Type     string                                     `json:"type"`
	External NotionPageExternalFileObjectExternalObject `json:"external"`
}

type NotionPageExternalFileCover struct {
	Type     string                                     `json:"type"`
	External NotionPageExternalFileObjectExternalObject `json:"external"`
}

type PokemonPageProperties struct {
	Name       NotionPageTitleProperty      `json:"Name"`
	No         NotionPageNumberObject       `json:"No"`
	Height     NotionPageNumberObject       `json:"Height"`
	Weight     NotionPageNumberObject       `json:"Weight"`
	HP         NotionPageNumberObject       `json:"HP"`
	Attack     NotionPageNumberObject       `json:"Attack"`
	Defense    NotionPageNumberObject       `json:"Defense"`
	SpAttack   NotionPageNumberObject       `json:"Sp. Attack"`
	SpDefense  NotionPageNumberObject       `json:"Sp. Defense"`
	Speed      NotionPageNumberObject       `json:"Speed"`
	Type       NotionPageMultiSelectObject  `json:"Type"`
	Sprite     NotionPageExternalFileObject `json:"Sprite"`
	Generation NotionPageSelectObject       `json:"Generation"`
}

type NotionPageTitleProperty struct {
	ID    string                 `json:"id"`
	Type  string                 `json:"type"`
	Title []NotionRichTextObject `json:"title"`
}

type NotionPageNumberObject struct {
	Number float64 `json:"number"`
}

type NotionPageMultiSelectObject struct {
	Type        string                         `json:"type"`
	MultiSelect []NotionPageSelectOptionObject `json:"multi_select"`
}

type NotionPageSelectObject struct {
	Type   string                       `json:"type"`
	Select NotionPageSelectOptionObject `json:"select"`
}

type NotionPageSelectOptionObject struct {
	Name string `json:"name"`
}

type NotionPageExternalFileObject struct {
	Type  string                                 `json:"type"`
	Files []NotionPageExternalFilePropertyObject `json:"files"`
}

type NotionPageExternalFilePropertyObject struct {
	Name     string                                     `json:"name"`
	Type     string                                     `json:"type"`
	External NotionPageExternalFileObjectExternalObject `json:"external"`
}

type NotionPageExternalFileObjectExternalObject struct {
	Url string `json:"url"`
}

type NotionPageBlockObject interface {
	BlockObject() string
	BlockType() string
}

type NotionBookmarkBlock struct {
	Object   string                  `json:"object"`
	Type     string                  `json:"type"`
	Bookmark NotionBookmarkBlockType `json:"bookmark"`
}

type NotionBookmarkBlockType struct {
	Url string `json:"url"`
}

type NotionCalloutBlock struct {
	Object  string                 `json:"object"`
	Type    string                 `json:"type"`
	Callout NotionCalloutBlockType `json:"callout"`
}

type NotionCalloutBlockType struct {
	RichText []NotionRichTextObject `json:"rich_text"`
	Icon     NotionEmojiObject      `json:"icon"`
	Color    string                 `json:"color"`
}

type NotionRichTextObject struct {
	Type        string                  `json:"type"`
	Text        NotionTextObject        `json:"text"`
	Annotations NotionAnnotationsObject `json:"annotations"`
	PlainText   string                  `json:"plain_text"`
}

type NotionTextObject struct {
	Content string `json:"content"`
}

type NotionAnnotationsObject struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

type NotionEmojiObject struct {
	Type  string `json:"type"`
	Emoji string `json:"emoji"`
}

func externalPokemonToInternalPokemon(external extract.Pokemon, databaseID DatabaseID) PokemonPage {
	output := PokemonPage{
		Parent: NotionPageParent{
			Type:       "database_id",
			DatabaseID: string(databaseID),
		},
		Icon: NotionPageExternalFileIcon{
			Type: "external",
			External: NotionPageExternalFileObjectExternalObject{
				Url: external.Sprite,
			},
		},
		Cover: NotionPageExternalFileCover{
			Type: "external",
			External: NotionPageExternalFileObjectExternalObject{
				Url: external.Artwork,
			},
		},
		Properties: PokemonPageProperties{
			Name: NotionPageTitleProperty{
				ID:   "title",
				Type: "title",
				Title: []NotionRichTextObject{
					{
						Type: "text",
						Text: NotionTextObject{
							Content: external.Name,
						},
						Annotations: NotionAnnotationsObject{
							Color: "default",
						},
						PlainText: external.Name,
					},
				},
			},
			No: NotionPageNumberObject{
				Number: float64(external.Number),
			},
			Height: NotionPageNumberObject{
				Number: float64(external.Height),
			},
			Weight: NotionPageNumberObject{
				Number: float64(external.Weight),
			},
			HP: NotionPageNumberObject{
				Number: float64(external.HP),
			},
			Attack: NotionPageNumberObject{
				Number: float64(external.Attack),
			},
			Defense: NotionPageNumberObject{
				Number: float64(external.Defense),
			},
			SpAttack: NotionPageNumberObject{
				Number: float64(external.SpAttack),
			},
			SpDefense: NotionPageNumberObject{
				Number: float64(external.SpDefense),
			},
			Speed: NotionPageNumberObject{
				Number: float64(external.Speed),
			},
			Type: NotionPageMultiSelectObject{
				Type:        "multi_select",
				MultiSelect: []NotionPageSelectOptionObject{},
			},
			Sprite: NotionPageExternalFileObject{
				Type: "files",
				Files: []NotionPageExternalFilePropertyObject{
					{
						Name: external.Sprite,
						Type: "external",
						External: NotionPageExternalFileObjectExternalObject{
							Url: external.Sprite,
						},
					},
				},
			},
			Generation: NotionPageSelectObject{
				Type: "select",
				Select: NotionPageSelectOptionObject{
					Name: external.Generation,
				},
			},
		},
		Children: []NotionPageBlockObject{},
	}

	for _, t := range external.Type {
		output.Properties.Type.MultiSelect = append(output.Properties.Type.MultiSelect, NotionPageSelectOptionObject{Name: t.Name})
	}

	for _, flavorText := range external.FlavorText {
		output.Children = append(output.Children, NotionCalloutBlock{
			Object: "block",
			Type:   "callout",
			Callout: NotionCalloutBlockType{
				RichText: []NotionRichTextObject{
					{
						Type: "text",
						Text: NotionTextObject{
							Content: fmt.Sprintf("%v\n", flavorText.Version),
						},
						Annotations: NotionAnnotationsObject{
							Bold:  true,
							Color: "default",
						},
					},
					{
						Type: "text",
						Text: NotionTextObject{
							Content: flavorText.Text,
						},
						Annotations: NotionAnnotationsObject{
							Color: "default",
						},
					},
				},
				Icon: NotionEmojiObject{
					Type:  "emoji",
					Emoji: "ℹ️",
				},
				Color: "default",
			},
		})
	}

	output.Children = append(output.Children, NotionBookmarkBlock{
		Object: "block",
		Type:   "bookmark",
		Bookmark: NotionBookmarkBlockType{
			Url: external.BulbapediaURL,
		},
	})

	return output
}

func (b NotionBookmarkBlock) BlockObject() string {
	return b.Object
}

func (b NotionBookmarkBlock) BlockType() string {
	return b.Type
}

func (b NotionCalloutBlock) BlockObject() string {
	return b.Object
}

func (b NotionCalloutBlock) BlockType() string {
	return b.Type
}
