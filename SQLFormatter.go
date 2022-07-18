package main

func sqlFormat(template string) ([]string, []int) {
	var (
		startPoints []int
		endPoints   []int
		returnArr   []string
		switchSpot  []int
	)
	for i, char := range template {
		if i+1 != len(template) {
			c := string(char)
			cp := string(template[i+1])
			if c == "$" && cp == "{" {
				startPoints = append(startPoints, i)
			}
		}
		if string(char) == "}" {
			endPoints = append(endPoints, i)
		}
	}
	if len(startPoints) != len(endPoints) {
		return nil, nil
	}
	returnArr = append(returnArr, string(template[0:startPoints[0]]))
	for i, startPoint := range startPoints {
		returnArr = append(returnArr, template[startPoint:endPoints[i]+1])
		switchSpot = append(switchSpot, len(returnArr)-1)
		if i+1 != len(startPoints) {
			returnArr = append(returnArr, template[endPoints[i]+1:startPoints[i+1]])
		} else {
			returnArr = append(returnArr, template[endPoints[i]+1:len(template)-1])
		}
	}
	return returnArr, switchSpot
}
