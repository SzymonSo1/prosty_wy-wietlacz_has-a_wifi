package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func main() {
	cmd := exec.Command("netsh", "wlan", "show", "profile")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Błąd podczas wykonywania polecenia:", err)
		return
	}

	result := string(output)

	re := regexp.MustCompile(`All User Profile\s*:\s*(.+)`)
	matches := re.FindAllStringSubmatch(result, -1)

	if len(matches) == 0 {
		fmt.Println("Nie znaleziono żadnych profili Wi-Fi.")
		return
	}

	var wifiInfo strings.Builder
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		profileName := strings.TrimSpace(match[1])
		wifiInfo.WriteString(fmt.Sprintf("\n=== Profil: %s ===\n", profileName))

		cmd := exec.Command("netsh", "wlan", "show", "profile", profileName, "key=clear")
		output, err := cmd.Output()
		if err != nil {
			wifiInfo.WriteString(fmt.Sprintf("Błąd podczas wykonywania polecenia dla profilu %s: %v\n", profileName, err))
			continue
		}

		keyContent := extractKeyContent(string(output))
		if keyContent != "" {
			wifiInfo.WriteString(fmt.Sprintf("Key Content (Hasło): %s\n", keyContent))
		} else {
			wifiInfo.WriteString("Key Content nie został znaleziony (hasło może być ukryte lub brak uprawnień).\n")
		}
	}

	fmt.Println(wifiInfo.String())

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Czy chcesz zapisać te informacje do pliku? (t/n): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "t" {
		file, err := os.Create("wifi_info.txt")
		if err != nil {
			fmt.Println("Błąd podczas tworzenia pliku:", err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(wifiInfo.String())
		if err != nil {
			fmt.Println("Błąd podczas zapisywania do pliku:", err)
			return
		}

		fmt.Println("Informacje zostały zapisane do pliku wifi_info.txt.")
	} else {
		fmt.Println("Informacje nie zostały zapisane.")
	}
	time.Sleep(10 * time.Second)
}

func extractKeyContent(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Key Content") {
			return strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}
	return ""
}