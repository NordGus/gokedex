package integrate

import (
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
	Category   NotionRichTextProperty       `json:"Category"`
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

type NotionRichTextProperty struct {
	Type     string                 `json:"type"`
	RichText []NotionRichTextObject `json:"rich_text"`
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

type NotionPageBlockObject interface {
	BlockObject() string
	BlockType() string
}

type NotionHeading3kBlock struct {
	Object   string                  `json:"object"`
	Type     string                  `json:"type"`
	Heading3 NotionHeading3BlockType `json:"heading_3"`
}

type NotionHeading3BlockType struct {
	RichText     []NotionRichTextObject `json:"rich_text"`
	Color        string                 `json:"color"`
	IsToggleable bool                   `json:"is_toggleable"`
}

func (b NotionHeading3kBlock) BlockObject() string {
	return b.Object
}

func (b NotionHeading3kBlock) BlockType() string {
	return b.Type
}

type NotionBookmarkBlock struct {
	Object   string                  `json:"object"`
	Type     string                  `json:"type"`
	Bookmark NotionBookmarkBlockType `json:"bookmark"`
}

type NotionBookmarkBlockType struct {
	Url string `json:"url"`
}

func (b NotionBookmarkBlock) BlockObject() string {
	return b.Object
}

func (b NotionBookmarkBlock) BlockType() string {
	return b.Type
}

type NotionTableBlock struct {
	Object string               `json:"object"`
	Type   string               `json:"type"`
	Table  NotionTableBlockType `json:"table"`
}

type NotionTableBlockType struct {
	TableWidth      float64               `json:"table_width"`
	HasColumnHeader bool                  `json:"has_column_header"`
	HasRowHeader    bool                  `json:"has_row_header"`
	Children        []NotionTableRowBlock `json:"children"`
}

type NotionTableRowBlock struct {
	Object   string                  `json:"object"`
	Type     string                  `json:"type"`
	TableRow NotionTableRowBlockType `json:"table_row"`
}

type NotionTableRowBlockType struct {
	Cells [][]NotionRichTextObject `json:"cells"`
}

func (b NotionTableBlock) BlockObject() string {
	return b.Object
}

func (b NotionTableBlock) BlockType() string {
	return b.Type
}

func externalPokemonToInternalPokemon(external extract.Pokemon, databaseID DatabaseID) PokemonPage {
	output := PokemonPage{}
	output.mapPageParent(external, databaseID)
	output.mapPageIcon(external)
	output.mapPageCover(external)
	output.Properties.mapPropertyName(external)
	output.Properties.mapPropertyCategory(external)
	output.Properties.No.mapNumericProperty(float64(external.Number))
	output.Properties.Height.mapNumericProperty(float64(external.Height))
	output.Properties.Weight.mapNumericProperty(float64(external.Weight))
	output.Properties.HP.mapNumericProperty(float64(external.HP))
	output.Properties.Attack.mapNumericProperty(float64(external.Attack))
	output.Properties.Defense.mapNumericProperty(float64(external.Defense))
	output.Properties.SpAttack.mapNumericProperty(float64(external.SpAttack))
	output.Properties.SpDefense.mapNumericProperty(float64(external.SpDefense))
	output.Properties.Speed.mapNumericProperty(float64(external.Speed))
	output.Properties.mapPropertyType(external)
	output.Properties.mapPropertySprite(external)
	output.Properties.mapPropertyGeneration(external)
	output.mapPageBlocks(external)

	return output
}

func (p *PokemonPage) mapPageParent(external extract.Pokemon, databaseId DatabaseID) {
	p.Parent = NotionPageParent{
		Type:       "database_id",
		DatabaseID: string(databaseId),
	}
}

func (p *PokemonPage) mapPageIcon(external extract.Pokemon) {
	p.Icon = NotionPageExternalFileIcon{
		Type: "external",
		External: NotionPageExternalFileObjectExternalObject{
			Url: external.Sprite,
		},
	}
}

func (p *PokemonPage) mapPageCover(external extract.Pokemon) {
	p.Cover = NotionPageExternalFileCover{
		Type: "external",
		External: NotionPageExternalFileObjectExternalObject{
			Url: external.Artwork,
		},
	}
}

func (p *PokemonPageProperties) mapPropertyName(external extract.Pokemon) {
	p.Name = NotionPageTitleProperty{
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
	}
}

func (p *PokemonPageProperties) mapPropertyCategory(external extract.Pokemon) {
	p.Category = NotionRichTextProperty{
		Type:     "rich_text",
		RichText: make([]NotionRichTextObject, len(external.Category)),
	}

	for i, category := range external.Category {
		p.Category.RichText[i] = NotionRichTextObject{
			Type: "text",
			Text: NotionTextObject{
				Content: category,
			},
			Annotations: NotionAnnotationsObject{
				Color: "default",
			},
			PlainText: category,
		}
	}
}

func (p *NotionPageNumberObject) mapNumericProperty(value float64) {
	p.Number = value
}

func (p *PokemonPageProperties) mapPropertyType(external extract.Pokemon) {
	p.Type = NotionPageMultiSelectObject{
		Type:        "multi_select",
		MultiSelect: make([]NotionPageSelectOptionObject, len(external.Type)),
	}

	for i, t := range external.Type {
		p.Type.MultiSelect[i] = NotionPageSelectOptionObject{Name: t.Name}
	}
}

func (p *PokemonPageProperties) mapPropertySprite(external extract.Pokemon) {
	p.Sprite = NotionPageExternalFileObject{
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
	}
}

func (p *PokemonPageProperties) mapPropertyGeneration(external extract.Pokemon) {
	p.Generation = NotionPageSelectObject{
		Type: "select",
		Select: NotionPageSelectOptionObject{
			Name: external.Generation,
		},
	}
}

func (p *PokemonPage) mapPageBlocks(external extract.Pokemon) {
	if len(external.FlavorText) > 0 {
		p.Children = make([]NotionPageBlockObject, 4)
		p.Children[0] = flavorTextTableHeadingBlock()
		p.Children[1] = flavorTextTableBlock(external.FlavorText)
		p.Children[2] = bulbapediaUrlBookmarkHeadingBlock()
		p.Children[3] = mapBulbapediaUrlBookmark(external)

		return
	}

	p.Children = make([]NotionPageBlockObject, 2)
	p.Children[0] = bulbapediaUrlBookmarkHeadingBlock()
	p.Children[1] = mapBulbapediaUrlBookmark(external)
}

func flavorTextTableHeadingBlock() NotionHeading3kBlock {
	return NotionHeading3kBlock{
		Object: "block",
		Type:   "heading_3",
		Heading3: NotionHeading3BlockType{
			RichText: []NotionRichTextObject{
				{
					Type: "text",
					Text: NotionTextObject{
						Content: "View This Pokémon's Pokédex Entries Through Generations:",
					},
					Annotations: NotionAnnotationsObject{
						Color: "default",
					},
					PlainText: "View This Pokémon's Pokédex Entries Through Generations:",
				},
			},
			Color: "default",
		},
	}
}

func flavorTextTableBlock(external []extract.PokemonFlavorText) NotionTableBlock {
	output := NotionTableBlock{
		Object: "block",
		Type:   "table",
		Table: NotionTableBlockType{
			TableWidth:      2,
			HasColumnHeader: true,
			HasRowHeader:    true,
			Children:        make([]NotionTableRowBlock, len(external)+1),
		},
	}

	output.Table.Children[0] = NotionTableRowBlock{
		Object: "block",
		Type:   "table_row",
		TableRow: NotionTableRowBlockType{
			Cells: [][]NotionRichTextObject{
				{
					{
						Type: "text",
						Text: NotionTextObject{
							Content: "Version",
						},
						Annotations: NotionAnnotationsObject{
							Bold:  true,
							Color: "default",
						},
						PlainText: "Version",
					},
				},
				{
					{
						Type: "text",
						Text: NotionTextObject{
							Content: "Flavor Text",
						},
						Annotations: NotionAnnotationsObject{
							Bold:  true,
							Color: "default",
						},
						PlainText: "Flavor Text",
					},
				},
			},
		},
	}

	for i, text := range external {
		output.Table.Children[i+1] = NotionTableRowBlock{
			Object: "block",
			Type:   "table_row",
			TableRow: NotionTableRowBlockType{
				Cells: [][]NotionRichTextObject{
					{
						{
							Type: "text",
							Text: NotionTextObject{
								Content: text.Version,
							},
							Annotations: NotionAnnotationsObject{
								Bold:  true,
								Color: "default",
							},
							PlainText: text.Version,
						},
					},
					{
						{
							Type: "text",
							Text: NotionTextObject{
								Content: text.Text,
							},
							Annotations: NotionAnnotationsObject{
								Color: "default",
							},
							PlainText: text.Text,
						},
					},
				},
			},
		}
	}

	return output
}

func bulbapediaUrlBookmarkHeadingBlock() NotionHeading3kBlock {
	return NotionHeading3kBlock{
		Object: "block",
		Type:   "heading_3",
		Heading3: NotionHeading3BlockType{
			RichText: []NotionRichTextObject{
				{
					Type: "text",
					Text: NotionTextObject{
						Content: "View This Pokémon's Entry on Bulbapedia:",
					},
					Annotations: NotionAnnotationsObject{
						Color: "default",
					},
					PlainText: "View This Pokémon's Entry on Bulbapedia:",
				},
			},
			Color: "default",
		},
	}
}

func mapBulbapediaUrlBookmark(external extract.Pokemon) NotionBookmarkBlock {
	return NotionBookmarkBlock{
		Object: "block",
		Type:   "bookmark",
		Bookmark: NotionBookmarkBlockType{
			Url: external.BulbapediaURL,
		},
	}
}
