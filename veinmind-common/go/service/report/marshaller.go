package report

import (
	"bytes"
	"encoding/json"
)

var (
	toLevel = map[Level]string{
		Low: "Low",
		Medium: "Medium",
		High: "High",
		Critical: "Critical",
	}

	fromLevel = map[string]Level{
		"Low": Low,
		"Medium": Medium,
		"High": High,
		"Critical": Critical,
	}

	toDetectType = map[DetectType]string{
		Image: "Image",
		Container: "Container",
	}

	fromDetectType = map[string]DetectType{
		"Image": Image,
		"Container": Container,
	}

	toEventType = map[EventType]string{
		Risk: "Risk",
		Invasion: "Invasion",
		Info: "Info",
	}

	fromEventType = map[string]EventType{
		"Risk": Risk,
		"Invasion": Invasion,
		"Info": Info,
	}

	toAlertType = map[AlertType]string{
		Vulnerability: "Vulnerability",
		MaliciousFile: "MaliciousFile",
		Backdoor: "Backdoor",
		Sensitive: "Sensitive",
		AbnormalHistory: "AbnormalHistory",
		Weakpass: "Weakpass",
		Asset: "Asset",
	}

	fromAlertType = map[string]AlertType{
		"Vulnerability": Vulnerability,
		"MaliciousFile": MaliciousFile,
		"Backdoor": Backdoor,
		"Sensitive": Sensitive,
		"AbnormalHistory": AbnormalHistory,
		"Weakpass": Weakpass,
		"Asset": Asset,
	}

	toWeakpassService = map[WeakpassService]string{
		SSH: "SSH",
	}

	fromWeakpassService = map[string]WeakpassService{
		"SSH": SSH,
	}
)

func (l Level) MarshalJSON()([]byte, error){
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toLevel[l])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (l *Level) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	switch v := i.(type) {
	case uint32:
		*l = Level(v)
	case float64:
		*l = Level(uint32(v))
	case string:
		*l = fromLevel[v]
	}

	return nil
}

func (d DetectType) MarshalJSON()([]byte, error){
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toDetectType[d])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (d *DetectType) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	switch i.(type) {
	case uint32:
		*d = i.(DetectType)
	case float64:
		*d = DetectType(i.(float64))
	case string:
		*d = fromDetectType[i.(string)]
	}

	return nil
}

func (e EventType) MarshalJSON()([]byte, error){
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toEventType[e])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (e *EventType) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	switch i.(type) {
	case uint32:
		*e = i.(EventType)
	case float64:
		*e = EventType(i.(float64))
	case string:
		*e = fromEventType[i.(string)]
	}

	return nil
}

func (a AlertType) MarshalJSON()([]byte, error){
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toAlertType[a])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (a *AlertType) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	switch i.(type) {
	case uint32:
		*a = i.(AlertType)
	case float64:
		*a = AlertType(i.(float64))
	case string:
		*a = fromAlertType[i.(string)]
	}

	return nil
}

func (w WeakpassService) MarshalJSON()([]byte, error){
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toWeakpassService[w])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (w *WeakpassService) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	switch i.(type) {
	case uint32:
		*w = i.(WeakpassService)
	case float64:
		*w = WeakpassService(i.(float64))
	case string:
		*w = fromWeakpassService[i.(string)]
	}

	return nil
}
