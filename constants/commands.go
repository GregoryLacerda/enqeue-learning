package constants

const (
	// Help Command
	CommandHelpMessage = `📚 **Comandos Disponíveis:**

				**!ping** - Testa se o bot está respondendo
				**!hello** / **!oi** - Recebe uma saudação
				**!calc <expressão>** - Calcula uma expressão matemática (ex: !calc 2 + 2)
				**!info** - Mostra informações sobre você e o sistema
				**!help** - Mostra esta mensagem de ajuda`

	// Ping Command
	PingMessage = "🏓 Pong!"

	// Hello Command
	HelloMessageTemplate = "👋 Olá, %s! Como posso ajudar?"

	// Info Command
	InfoMessageTemplate = `ℹ️ **Informações do Sistema:**

👤 **Usuário:** %s
🆔 **User ID:** %s
📝 **Comando:** %s
📅 **Canal ID:** %s
🏢 **Guild ID:** %s
⏰ **Timestamp:** %s`

	// Unknown Command
	UnknownCommandTemplate = "❓ Comando desconhecido: `%s`\n\nUse `!help` para ver os comandos disponíveis."

	// Calc Command
	CalcUsageMessage        = "❌ Uso: `!calc <expressão>`\n**Exemplo:** !calc 2 + 2"
	CalcErrorTemplate       = "❌ Erro ao calcular: %s"
	CalcResultTemplate      = "🧮 **Resultado:**\n`%s = %.2f`"
	CalcInvalidFormat       = "formato esperado: número operador número (ex: 2 + 2)"
	CalcInvalidFirstNumber  = "primeiro número inválido: %s"
	CalcInvalidSecondNumber = "segundo número inválido: %s"
	CalcDivisionByZero      = "divisão por zero"
	CalcInvalidOperator     = "operador inválido: %s (use: +, -, *, /, ^)"
)
