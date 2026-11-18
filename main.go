package main

import (
 "bufio"
 "fmt"
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "syscall"
 "unsafe"
)

const title = "Instalador GPEDIT v1.0"

func main() {
 fmt.Println("============================================================")
 fmt.Println("  " + title)
 fmt.Println("  Instala pacotes nativos de Politica de Grupo do Windows")
 fmt.Println("============================================================")
 fmt.Println()
 fmt.Println("AVISO: Windows Home NAO inclui gpedit.msc como recurso")
 fmt.Println("oficialmente suportado. Este procedimento NAO transforma")
 fmt.Println("o Windows Home em Pro e NAO garante o funcionamento das")
 fmt.Println("politicas de grupo nessa edicao, mesmo sem erros no DISM.")
 fmt.Println()
 if !isAdmin() {
  fmt.Println("Solicitando permissao de administrador...")
  if err := elevate(); err != nil { fmt.Println("Nao foi possivel elevar: ",err); pause() }
  return
 }
 defer pause()
 root := os.Getenv("SystemRoot")
 if root == "" { root = `C:\Windows` }
 servicing := filepath.Join(root,"servicing","Packages")
 kinds := []string{"Microsoft-Windows-GroupPolicy-ClientTools-Package~*.mum", "Microsoft-Windows-GroupPolicy-ClientExtensions-Package~*.mum"}
 var files []string
 for _, kind := range kinds {
  found, err := filepath.Glob(filepath.Join(servicing,kind))
  if err != nil || len(found)==0 { fmt.Printf("[ERRO] Nenhum pacote encontrado: %s\n",kind); fmt.Println("Nao e possivel instalar: os pacotes estao ausentes nesta imagem do Windows."); return }
  fmt.Printf("[OK] %d pacote(s) localizado(s): %s\n",len(found), strings.Split(kind,"~")[0])
  files=append(files,found...)
 }
 fmt.Printf("\nPacotes a processar: %d\n",len(files))
 fmt.Println("Usando DISM /Online /NoRestart /Add-Package. Nenhum arquivo externo sera baixado.")
 fmt.Print("\nContinuar? (S/N): ")
 in:=bufio.NewReader(os.Stdin)
 answer,_:=in.ReadString('\n')
 if strings.ToLower(strings.TrimSpace(answer))!="s" { fmt.Println("Operacao cancelada."); return }
 dism:=filepath.Join(root,"System32","dism.exe")
 var failed, reboot int
 for i, file := range files {
  fmt.Printf("\n[%d/%d] Instalando %s\n", i+1, len(files), filepath.Base(file))
  cmd:=exec.Command(dism,"/Online","/NoRestart","/Add-Package:"+file)
  cmd.Stdout=os.Stdout;cmd.Stderr=os.Stderr
  err:=cmd.Run()
  if err!=nil {
   if exit,ok:=err.(*exec.ExitError);ok && exit.ExitCode()==3010 {reboot++;fmt.Println("[OK] Instalado. Reinicializacao necessaria.");continue}
   failed++;fmt.Printf("[ERRO] DISM falhou: %v\n",err)
  } else {fmt.Println("[OK] Comando finalizado.")}
 }
 fmt.Println("\n============================================================")
 if failed>0 {fmt.Printf("Instalacao com falhas: %d de %d pacote(s).\n",failed,len(files))} else {fmt.Println("Processamento dos pacotes concluido sem erros DISM.")}
 if reboot>0 {fmt.Println("Reinicie o Windows antes de testar o gpedit.msc.")}
 fmt.Println("Abra Win+R e digite gpedit.msc para verificar o funcionamento.")
 fmt.Println("Nota: o Windows Home nao oferece suporte oficial ao GPEDIT.")
}

func isAdmin() bool {
 cmd:=exec.Command("fltmc.exe")
 return cmd.Run()==nil
}
func elevate() error {
 shell:=syscall.NewLazyDLL("shell32.dll").NewProc("ShellExecuteW")
 exe,err:=os.Executable(); if err!=nil{return err}
 str:=func(s string)*uint16 {p,_:=syscall.UTF16PtrFromString(s);return p}
 r,_,e:=shell.Call(0,uintptr(unsafe.Pointer(str("runas"))),uintptr(unsafe.Pointer(str(exe))),0,0,1)
 if r<=32 {return fmt.Errorf("Windows retornou %d (%v)",r,e)}
 return nil
}
func pause(){fmt.Print("\nPressione ENTER para fechar...");bufio.NewReader(os.Stdin).ReadString('\n')}
