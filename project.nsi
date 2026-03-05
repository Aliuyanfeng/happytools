Unicode true

!include "wails_tools.nsh"

VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

ManifestDPIAware true

!include "MUI.nsh"

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"

!define MUI_WELCOMEPAGE_TITLE "欢迎使用 HappyTools 安装向导"
!define MUI_WELCOMEPAGE_TEXT "这将安装 HappyTools 到您的计算机。$\r$\n$\r$\nHappyTools 是一款提升幸福感的工具集，包含系统监控、待办管理、VT扫描等功能。$\r$\n$\r$\n建议您在继续之前关闭所有其他应用程序。"

!define MUI_FINISHPAGE_TITLE "HappyTools 安装完成"
!define MUI_FINISHPAGE_TEXT "HappyTools 已成功安装到您的计算机。$\r$\n$\r$\n点击完成退出安装向导。"
!define MUI_FINISHPAGE_RUN "$INSTDIR\${PRODUCT_EXECUTABLE}"
!define MUI_FINISHPAGE_RUN_TEXT "立即运行 HappyTools"
!define MUI_FINISHPAGE_SHOWREADME ""
!define MUI_FINISHPAGE_SHOWREADME_NOTCHECKED
!define MUI_FINISHPAGE_SHOWREADME_TEXT "查看 README 文档"
!define MUI_FINISHPAGE_SHOWREADME_FUNCTION "ShowReadme"

!define MUI_DIRECTORYPAGE_TEXT_TOP "选择 HappyTools 的安装位置"

!define MUI_COMPONENTSPAGE_SMALLDESC
!define MUI_FINISHPAGE_NOAUTOCLOSE

!define MUI_ABORTWARNING
!define MUI_ABORTWARNING_TEXT "您确定要退出 HappyTools 安装程序吗？"

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "SimpChinese"

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\..\bin\${INFO_PROJECTNAME}-${ARCH}-installer.exe"
InstallDir "$PROGRAMFILES64\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}"
ShowInstDetails show

Function .onInit
   !insertmacro wails.checkArchitecture
FunctionEnd

Function ShowReadme
    ExecShell "open" "https://github.com/Aliuyanfeng/happytools"
FunctionEnd

Section "!HappyTools (必选)" SecMain
    SectionIn RO
    !insertmacro wails.setShellContext
    !insertmacro wails.webview2runtime
    SetOutPath $INSTDIR
    !insertmacro wails.files
    !insertmacro wails.associateFiles
    !insertmacro wails.writeUninstaller
SectionEnd

Section "创建桌面快捷方式" SecDesktop
    CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
SectionEnd

Section "创建开始菜单快捷方式" SecStartMenu
    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
SectionEnd

LangString DESC_SecMain ${LANG_SIMPCHINESE} "安装 HappyTools 核心文件（必选）"
LangString DESC_SecDesktop ${LANG_SIMPCHINESE} "在桌面创建快捷方式"
LangString DESC_SecStartMenu ${LANG_SIMPCHINESE} "在开始菜单创建快捷方式"

!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
  !insertmacro MUI_DESCRIPTION_TEXT ${SecMain} $(DESC_SecMain)
  !insertmacro MUI_DESCRIPTION_TEXT ${SecDesktop} $(DESC_SecDesktop)
  !insertmacro MUI_DESCRIPTION_TEXT ${SecStartMenu} $(DESC_SecStartMenu)
!insertmacro MUI_FUNCTION_DESCRIPTION_END

Section "uninstall"
    !insertmacro wails.setShellContext
    RMDir /r "$AppData\${PRODUCT_EXECUTABLE}"
    RMDir /r $INSTDIR
    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"
    !insertmacro wails.unassociateFiles
    !insertmacro wails.deleteUninstaller
SectionEnd
