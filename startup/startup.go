package startup

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/config"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/defaults"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/script"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/statistic"
)

func InitLogs() {
	file, err := os.OpenFile("./vmupdater.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}

	logrus.SetOutput(file)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true,
		TimestampFormat: "2006/01/02 15:04:05"})
}

func RecoverMain() {
	if r := recover(); r != nil {
		logrus.Errorln("panic: ", r)
		fmt.Println("panic: ", r)
	}
}

func Prepare() {
	InitLogs()
	// Создать папку в дефолтной директории пользователя если отсутствует
	// На винде это Documents/updater, на линуксе ~/.config/updater
	// Дира .*/updater
	// В ней лежит config.json предназначенный для хранения всей общей информации, которую может настроить пользователь
	config.ReadDefaultFile()

	if config.Empty() {
		defaults.FillSupported()
	}
	// Грузим скрипты для обновления сервисов из .*/updater/scripts/
	script.Init()
	// Грузим данные по статистике (для отображения процентов)
	statistic.ReadFromDefaultFile()
}
