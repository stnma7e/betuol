package character

import (
	"fmt"

	"smig/component"
)

type AiComputer func(id component.GOiD, chars *CharacterManager)

func (cm *CharacterManager) RegisterComputer(name string, aic AiComputer) {
	if cm.aiFunctionName == nil {
		cm.aiFunctionName = make(map[string]AiComputer)
	}
	cm.aiFunctionName[name] = aic
}
func (cm *CharacterManager) GetComputer(name string) (AiComputer, error) {
	aic, ok := cm.aiFunctionName[name]
	if !ok {
		return func(component.GOiD, *CharacterManager){}, fmt.Errorf("unregistered AiComputer name: %s", name)
	}
	return aic, nil
}

func (cm *CharacterManager) RunAi(id component.GOiD) {
	cm.aiList[id](id, cm)
}

/*****************************************
*
* AI Functions
*
*****************************************/

func PassiveComputer(id component.GOiD, chars *CharacterManager) {
	fmt.Println("passive")
}

func PlayerComputer(id component.GOiD, chars *CharacterManager) {
	ca := chars.GetCharacterAttributes(id)
	loc :=  chars.Scene.GetObjectLocation(id)
	fmt.Print(ca.Attributes[HEALTH], ca.Attributes[MANA], loc[:2], " --> ")

	var command string
	fmt.Scan(&command)
	ParsePlayerCommand(command, id, chars)
}

func MerchantComputer(id component.GOiD, chars *CharacterManager) {
	fmt.Println("merchanting")
}