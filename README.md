# Instalador GPEDIT v1.0

Utilitario para Windows x64 baseado em comandos DISM. Procura e processa pacotes locais `Microsoft-Windows-GroupPolicy-ClientTools-Package` e `Microsoft-Windows-GroupPolicy-ClientExtensions-Package`.

## ATENCAO: Windows Home

**O Windows Home NAO inclui o Editor de Politica de Grupo Local (`gpedit.msc`) como recurso oficialmente disponivel/suportado pela Microsoft.**

Este utilitario **NAO transforma o Windows Home em Pro**, **NAO garante compatibilidade com o Home** e **NAO garante que as politicas exibidas no editor sejam aplicadas pelo Windows Home**. Mesmo que o DISM conclua os comandos sem erro, o funcionamento completo do editor e de suas politicas nao e garantido.

Para administracao de Politica de Grupo com suporte oficial, utilize uma edicao compativel do Windows, como Pro, Enterprise ou Education.

## Funcionamento

- Solicita permissao administrativa via UAC.
- Busca pacotes `.mum` locais em `%SystemRoot%\servicing\Packages`.
- Exibe os pacotes encontrados e solicita confirmacao antes das alteracoes.
- Processa cada pacote com `DISM /Online /NoRestart /Add-Package`.
- Exibe erros do DISM, quando houver.
- Nao baixa pacotes, nao modifica diretamente o Registro e nao oferece suporte oficial ao Windows Home.

## Como utilizar

1. Baixe a versao publicada em **Releases**.
2. Execute `Instalar_GPEDIT_v1.0.exe` como administrador e confirme a solicitacao UAC.
3. Revise os pacotes encontrados e digite `S` para prosseguir.
4. Aguarde o processamento e observe os resultados. Reinicie o Windows, se solicitado.
5. Execute `gpedit.msc` para conferir o comportamento no computador.

**Recomendacao:** crie um ponto de restauracao e utilize o procedimento somente em maquinas autorizadas. Uma execucao do DISM sem erros nao valida todas as politicas.

## Compilacao

Requer Go 1.22 ou superior:

```powershell
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
$env:CGO_ENABLED = '0'
go build -trimpath -ldflags='-s -w' -o Instalar_GPEDIT_v1.0.exe .
```

O workflow `.github/workflows/release.yml` compila para Windows x64, calcula SHA-256 e publica o executavel na Release `v1.0.0`.

O executavel nao tem assinatura digital. Verifique o arquivo antes de distribuir. Teste inicial realizado em Windows pelo solicitante; nao ha certificacao de funcionamento em todas as versoes e edicoes do Windows.

Repositorio de distribuicao independente do aplicativo AjudaSuporte.
