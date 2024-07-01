package gui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"github.com/sirupsen/logrus"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/config"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/updater"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func UpdateAction(button *widget.Button, name string, u *updater.Updater) {
	var err error

	// Через эту штуку updater с интерфейсом связан
	progress := make(chan util.Progress)
	defer close(progress)

	go func() {
		for val := range progress {
			// А все что он пишет - будем рисовать в интерфейсе
			button.SetText(val.String())
		}

		// progress is closed
		if err == nil {
			button.SetText(fmt.Sprintf("%s updated", name))
			button.SetIcon(theme.ConfirmIcon())
		}

		u.SetBusy(false)
	}()

	// Вот тут обновление и происходит
	err = u.Update(progress)
	if err != nil {
		progress <- util.Progress{Error: err}
	}
}

const DefaultRect = 1000

func Application() {
	// gui инициализируем
	a := app.New()
	w = a.NewWindow("Updater")
	w.Resize(fyne.NewSize(DefaultRect, DefaultRect))

	// Получаем список поддерживаемых сервисов из config.json
	service2Supported := config.GetSupported()

	// Инициализируем кнопочки
	localButtons := LocalButtons(service2Supported)
	remoteButtons := RemoteButtons(service2Supported)
	// inplaceButtons := InplaceButtons(service2Supported)

	// Запихиваем их в сеточки для отображения
	// TODO починить вариадики
	// TODO Рассчитывать количество столбцов из ширины экрана
	local := container.New(layout.NewAdaptiveGridLayout(2)) //nolint:gomnd
	for _, b := range localButtons {
		local.Add(b)
	}

	remote := container.New(layout.NewAdaptiveGridLayout(2)) //nolint:gomnd
	for _, b := range remoteButtons {
		remote.Add(b)
	}

	// inplace := container.New(layout.NewAdaptiveGridLayout(2)) //nolint:gomnd
	// for _, b := range inplaceButtons {
	// 	inplace.Add(b)
	// }

	// Добавляем пару кнопок конфигурации
	// TODO утащить куда нибудь
	configButton := ConfigButton()
	clearDirsButton := ClearDirsButton()

	local.Add(configButton)
	local.Add(clearDirsButton)
	remote.Add(configButton)
	remote.Add(clearDirsButton)
	// inplace.Add(configButton)
	// inplace.Add(clearDirsButton)

	// Вкладочки, в одной обновление со своей машины, в другой напрямую на сервере
	tabs := container.NewAppTabs(
		container.NewTabItem("Local", local),
		container.NewTabItem("Container", remote),
		// container.NewTabItem("Inplace", inplace),
	)

	w.SetContent(tabs)
	w.Show()
	a.Run()
	w.Close()
	a.Driver().Quit()
}

// Настройки локального обновления.
func LocalConfigWindow(name string) {
	dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
		if uri != nil {
			// Диалог вылезает с выбором диры. Мы ее пишем в настройки сервиса, чтобы оттуда файлы грузить
			config.Set(name, "dir", uri.Path())
		}
	}, w)
}

const defaultMinSize = 200

// Настройки обновления сразу на сервере.
func RemoteConfigWindow(name string, confirmed chan bool) {
	// Данные для выпадашки с выбором веток
	selects, _ := config.GetSelects(name + "_branch")
	branch := widget.NewSelectEntry(selects)
	branch.SetText("master")
	// Без бранчи из этого окна не уйти ╮(-_-)╭
	branch.Validator = validation.NewRegexp(`[^\s]`, "Branch shoud be specified")

	// В окно по умолчанию подставляем последнюю выбранную ветку
	curBranch, ok := config.Get(name, "branch")
	if ok {
		branch.SetText(curBranch)
	}

	// Собственно само окошечко со всей инфой
	items := []*widget.FormItem{{Text: "Branch", Widget: branch}}
	d := dialog.NewForm("Choose branch", "Update", "Dismiss", items, func(b bool) {
		confirmed <- b
		if !b {
			return
		}
		// Если ок нажали - продолжаем
		br := branch.Text
		// Ставим текущую бранчу для обновления
		config.Set(name, "branch", br)
		logrus.Debugln("Choosed branch", br)

		// Добавляем её в хранение выпадашек
		config.AddSelect(name+"_branch", br)

		// Обновляем выпадашки, чтобы при повторном открытии эта бранча там была
		selects, _ = config.GetSelects(name + "_branch")
		branch.SetOptions(selects)
	}, w)
	d.Resize(d.MinSize().AddWidthHeight(defaultMinSize, 0))
	d.Show()
}

// Кнопка обновления из локальной сборки.
func LocalButton(name string) *widget.Button {
	button := widget.NewButton(name, nil)
	button.OnTapped = func() {
		// Если папка со сборкой не настроена
		_, ok := config.Get(name, "dir")
		if !ok {
			// То настроим
			LocalConfigWindow(name)

			return
		}

		u := updater.Get(name, util.UpdateLocal, config.Args(name), config.Args("ssh"))
		if u.IsBusy() {
			// Если updater уже работает - нечего его в параллель загружать
			fmt.Println("busy updater, do nothing")

			return
		}

		// Иначе приступаем к обновлению
		u.SetBusy(true)
		button.SetText("0%")
		button.SetIcon(theme.DownloadIcon())

		go UpdateAction(button, name, u)
	}

	return button
}

const defaultTimeout = 5 * time.Minute

// Кнопка обновления напрямую на сервере.
func RemoteButton(name string) *widget.Button {
	button := widget.NewButton(name, nil)
	button.OnTapped = func() {
		// Для связи с вылезающим окошком выбора бранчи
		confirmed := make(chan bool)

		u := updater.Get(name, util.UpdateRemote, config.Args(name), config.Args("ssh"))
		if u.IsBusy() {
			// Если updater уже работает - нечего его в параллель загружать
			fmt.Println("busy updater, do nothing")

			return
		}
		// Запускаем диалоговое окно.
		RemoteConfigWindow(name, confirmed)

		go func() {
			select {
			case b := <-confirmed:
				// Если в нем нажали отклонить - обновляться не стоит
				if !b {
					return
				}

				u.SetBusy(true)
				button.SetText("0%")
				button.SetIcon(theme.DownloadIcon())
				// Эта функция уже горутина, не блокирующая интерфейс. Поэтому можно синхронно запустить
				UpdateAction(button, name, u)
			case <-time.After(defaultTimeout):
				return
			}
		}()
	}

	return button
}

// Кнопка обновления из локальной сборки.
// func InplaceButton(name string) *widget.Button {
// 	button := widget.NewButton(name, nil)
// 	button.OnTapped = func() {
// 		u := updater.Get(name, util.UpdateInPlace, config.Args(name))
// 		if u.IsBusy() {
// 			// Если updater уже работает - нечего его в параллель загружать
// 			fmt.Println("busy updater, do nothing")
//
// 			return
// 		}
//
// 		// Иначе приступаем к обновлению
// 		u.SetBusy(true)
// 		button.SetText("0%")
// 		button.SetIcon(theme.DownloadIcon())
//
// 		go UpdateAction(button, name, u)
// 	}
//
// 	return button
// }

func ConfigButton() *widget.Button {
	return widget.NewButton("SSH settings", func() {
		d := SSHConfigDialog()
		d.Resize(d.MinSize().AddWidthHeight(defaultMinSize, 0))
		d.Show()
	})
}

func AddrSelectEntry() *widget.SelectEntry {
	selects, _ := config.GetSelects("addr")
	addr := widget.NewSelectEntry(selects)
	addr.Validator = validation.NewRegexp(`[^\s]`, "Address shoud be specified")

	addrT, ok := config.Get("ssh", "addr")
	if ok {
		addr.Text = addrT
	}

	return addr
}

func SSHEntry(opt string) (*widget.Entry, string) {
	key := widget.NewEntry()

	keyT, ok := config.Get("ssh", opt)
	if ok {
		key.Text = keyT
	}

	return key, keyT
}

func KeyBut(key *widget.Entry) *widget.Button {
	return widget.NewButtonWithIcon("", theme.DocumentIcon(), func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri != nil {
				path := uri.URI().Path()
				config.Set("ssh", "key", path)
				key.Text = path
			}
		}, w)
	})
}

func PassEntry() *widget.Entry {
	pass := widget.NewPasswordEntry()

	passT, ok := config.Get("ssh", "pass")
	if ok {
		pass.Text = passT
	}

	return pass
}

func SSHConfigDialog() dialog.Dialog {
	key, keyT := SSHEntry("key")
	keyBut := KeyBut(key)
	user, _ := SSHEntry("user")
	pass := PassEntry()

	addr := AddrSelectEntry()
	items := []*widget.FormItem{ // we can specify items in the constructor
		{Text: "Address", Widget: addr},
		{Text: "User", Widget: user},
		{Text: "Password", Widget: pass},
		{Text: "Key", Widget: key},
		{Text: "Choose key", Widget: keyBut},
	}

	return dialog.NewForm("SSH settings", "Confirm", "Dismiss", items, func(b bool) {
		if !b {
			return
		}
		address := addr.Text
		config.Set("ssh", "addr", address)
		config.AddSelect("addr", address)
		s, _ := config.GetSelects("addr")
		addr.SetOptions(s)
		if user.Text != "" {
			config.Set("ssh", "user", user.Text)
		} else {
			config.Unset("ssh", "user")
		}
		if keyT != "" {
			config.Set("ssh", "key", keyT)
		} else {
			config.Unset("ssh", "key")
		}
		if pass.Text != "" {
			config.Set("ssh", "pass", pass.Text)
		} else {
			config.Unset("ssh", "pass")
		}
	}, w)
}

func ClearDirsButton() *widget.Button {
	return widget.NewButton("Clear directories", func() { config.ClearDirs() })
}

func LocalButtons(services map[string]string) [](*widget.Button) {
	return GetButtons(services, "1357", LocalButton)
}
func RemoteButtons(services map[string]string) (l [](*widget.Button)) {
	return GetButtons(services, "2367", RemoteButton)
}

// func InplaceButtons(services map[string]string) (l [](*widget.Button)) {
// 	return GetButtons(services, "4567", InplaceButton)
// }

func GetButtons(services map[string]string, supported string, but func(string) *widget.Button) (l [](*widget.Button)) {
	var names []string

	for k, v := range services {
		// 0 - Not supported
		// 1 - Local
		// 2 - Remote
		// 3 - Local | Remote
		// 4 - Inplace
		// 5 - Local | Inplace
		// 6 - Remote | Inplace
		// 7 - Local | Remote | Inplace
		if strings.Contains(supported, v) {
			names = append(names, k)
		}
	}

	sort.Strings(names)

	for _, k := range names {
		l = append(l, but(k))
	}

	return
}

// Требуется в любой всплывашке/диалоге и т.д.
// Пусть будет глобальным, я его все равно не меняю напрямую, только пихаю везде.
var w fyne.Window //nolint:gochecknoglobals
