package options

import (
	"financial-cli/internal/application"
	"financial-cli/internal/cli/view"
	"financial-cli/internal/config/utils"
	"financial-cli/internal/domain/category"
	"financial-cli/internal/domain/ofx"
	"financial-cli/internal/domain/register"
	"financial-cli/internal/domain/transactions"

	"fmt"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func ImportCommand(db *gorm.DB) *cli.Command {
	return &cli.Command{
		Name:  "import",
		Usage: "Importa um arquivo OFX",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Caminho para o arquivo OFX",
				Value:   "./import.ofx",
			},
		},
		Action: func(c *cli.Context) error {
			view.RunIfNotDebug(view.ClearScreen)
			path := c.Path("path")

			ofxData, err := ofx.ParseOFXFile(path)
			if err != nil {
				return fmt.Errorf("erro ao ler arquivo OFX: %w", err)
			}

			categoryRepo := category.NewCategoryRepository(db)
			registerRepo := register.NewRegisterRepository(db)
			transactionRepo := transactions.NewTransactionRepository(db)

			app := application.NewApplication(categoryRepo, transactionRepo, registerRepo)

			categories, err := categoryRepo.GetAll()
			if err != nil {
				return fmt.Errorf("erro ao buscar categorias: %w", err)
			}

			reg := register.Register{
				OrgID:        ofxData.Signon.FID,
				Account:      ofxData.BankResponse.AccountID,
				StartDate:    utils.ParseOFXDate(ofxData.BankResponse.StartDate),
				EndDate:      utils.ParseOFXDate(ofxData.BankResponse.EndDate),
				Organization: ofxData.Signon.Org,
				Amount:       ofxData.BankResponse.Balance,
			}

			regID, err := app.OFXImporter.CreateRegister(reg)
			if err != nil {
				return fmt.Errorf("erro ao creiar Registro: %w", err)
			}

			for _, item := range ofxData.BankResponse.Transactions {
				tx := transactions.Transaction{
					Description:     item.Description,
					Date:            utils.ParseOFXDate(item.Date),
					Value:           item.Amount,
					TransactionType: item.Type,
					TransactionID:   item.ID,
					RegisterID:      regID,
				}
				existing, _ := app.OFXImporter.SearchDuplicateTransaction(tx.Value, tx.Date, tx.Description, tx.TransactionID)
				if existing != nil {
					fmt.Printf("Transação encontrada. Pulando transação - %s.\n", tx.Description)
					view.RunIfNotDebug(view.ClearScreen)
					continue
				}

				catID, err := view.PromptCategory(categories, item.Description, item.Amount, item.Date)
				if err != nil {
					fmt.Println("Entrada inválida. Pulando transação.")
					continue
				}
				tx.CategoryID = catID

				note, err := view.PromptNote()

				if err != nil {
					fmt.Println("Entrada inválida. Pulando transação.")
					continue
				}
				tx.Note = note

				if err := app.OFXImporter.ImportTransaction(tx); err != nil {
					return fmt.Errorf("erro ao importar transações: %w", err)
				}
			}
			view.RunIfNotDebug(view.ClearScreen)
			fmt.Println("Importação concluída com sucesso!")
			return nil
		},
	}
}
