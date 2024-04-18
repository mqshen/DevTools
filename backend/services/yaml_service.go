package services

import (
	"context"
	"devtools/backend/types"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// NewLine 换行符
var NewLine = "\n"

// SignEqual 等号连接符
var SignEqual = "="

// NewLineDom yaml的value换行符
var YamlNewLineDom = "|\n"

// IndentBlanks 缩进空格
var IndentBlanks = "  "

// SignSemicolon 分号连接符
var SignSemicolon = ":"

// ArrayBlanks 数组缩进
var ArrayBlanks = "- "

var rangePattern = regexp.MustCompile("^(.*)\\[(\\d*)\\]$")

type yamlService struct {
	ctx context.Context
}

var yamlConvertor *yamlService
var onceYamlConvertor sync.Once

func YamlConvertor() *yamlService {
	if yamlConvertor == nil {
		onceYamlConvertor.Do(func() {
			yamlConvertor = &yamlService{}
		})
	}
	return yamlConvertor
}

func (p *yamlService) Start(ctx context.Context) {
	p.ctx = ctx
}

func (p *yamlService) ConvertFromJSON(content string) (resp types.JSResp) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	yamlContent, err := yaml.Marshal(data)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	jsonContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Data = map[string]any{
		"yamlContent": string(yamlContent),
		"jsonContent": string(jsonContent),
	}
	resp.Success = true
	return
}

func (p *yamlService) ConvertToJSON(content string) (resp types.JSResp) {
	var data map[string]interface{}
	err := yaml.Unmarshal([]byte(content), &data)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	yamlContent, err := yaml.Marshal(data)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	jsonContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Data = map[string]any{
		"yamlContent": string(yamlContent),
		"jsonContent": string(jsonContent),
	}
	resp.Success = true
	return
}

func (p *yamlService) ConvertFromProperties(properties string) (resp types.JSResp) {
	var yamlLineList []string
	var yamlNodes []types.YamlNode
	propertiesLineWordList := GetPropertiesItemLineList(properties)
	for _, line := range propertiesLineWordList {
		line = strings.TrimSpace(line)
		if line != "" {
			// 注释数据不要
			if strings.HasPrefix(line, "#") {
				continue
			}

			index := strings.Index(line, SignEqual)
			if index > -1 {
				key := line[:index]
				value := line[index+1:]

				if strings.Contains(value, "\n") {
					value = YamlNewLineDom + value
				}

				lineWordList := strings.Split(key, ".")
				lineWordList, yamlNodes = p.wordToNode(lineWordList, yamlNodes, nil, false, -1, appendSpaceForArrayValue(value))
			}
		}
	}
	yamlLineList = formatPropertiesToYaml(yamlLineList, yamlNodes, false, "")
	resp.Success = true
	resp.Data = strings.Join(yamlLineList, "\n") + "\n"
	return
}

func (p *yamlService) ConvertToProperties(content string) (resp types.JSResp) {
	// yaml 到 map
	dataMap, err := p.yamlToMap(content)
	if err != nil {
		runtime.LogErrorf(p.ctx, "YamlToPropertiesStr error: %v, content: %v", err, content)
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	data, err := p.mapToProperties(dataMap)
	if err != nil {
		runtime.LogErrorf(p.ctx, "YamlToPropertiesStr error: %v, content: %v", err, content)
		resp.Msg = err.Error()
		return
	}
	resp.Data = data
	return
}

func formatPropertiesToYaml(yamlLineList []string, yamlNodes []types.YamlNode, lastNodeArrayFlag bool, blanks string) []string {
	var beforeNodeIndex = -1
	var equalSign string

	for _, yamlNode := range yamlNodes {
		value := yamlNode.Value

		equalSign = SignSemicolon
		if "" != value {
			equalSign = SignSemicolon + " "
		}

		yamlNode.ResortValue()

		name := yamlNode.Name
		if lastNodeArrayFlag {
			if "" == name {
				yamlLineList = append(yamlLineList, blanks+ArrayBlanks+stringValueWrap(value))
			} else {
				if -1 != beforeNodeIndex && beforeNodeIndex == yamlNode.LastNodeIndex {
					yamlLineList = append(yamlLineList, blanks+IndentBlanks+name+equalSign+stringValueWrap(value))
				} else {
					yamlLineList = append(yamlLineList, blanks+ArrayBlanks+name+equalSign+stringValueWrap(value))
				}
			}
			beforeNodeIndex = yamlNode.LastNodeIndex
		} else {
			yamlLineList = append(yamlLineList, blanks+name+equalSign+stringValueWrap(value))
		}

		if yamlNode.ArrayFlag {
			if lastNodeArrayFlag {
				yamlLineList = formatPropertiesToYaml(yamlLineList, yamlNode.ValueList, true, IndentBlanks+IndentBlanks+blanks)
			} else {
				yamlLineList = formatPropertiesToYaml(yamlLineList, yamlNode.ValueList, true, IndentBlanks+blanks)
			}
		} else {
			if lastNodeArrayFlag {
				yamlLineList = formatPropertiesToYaml(yamlLineList, yamlNode.Children, false, IndentBlanks+IndentBlanks+blanks)
			} else {
				yamlLineList = formatPropertiesToYaml(yamlLineList, yamlNode.Children, false, IndentBlanks+blanks)
			}
		}
	}
	return yamlLineList
}

func (p *yamlService) yamlToMap(contentOfYaml string) (map[string]interface{}, error) {
	resultMap := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(contentOfYaml), &resultMap)
	if err != nil {
		runtime.LogErrorf(p.ctx, "YamlToMap, error: %v, content: %v", err, contentOfYaml)
		return nil, err
	}

	return resultMap, nil
}

// 进行深层嵌套的map数据处理
func (p *yamlService) mapToProperties(dataMap map[string]interface{}) (string, error) {
	var propertyStrList []string
	for key, value := range dataMap {
		if value == nil {
			continue
		}
		valueKind := reflect.TypeOf(value).Kind()
		switch valueKind {
		case reflect.Map:
			{
				propertyStrList = doMapToProperties(propertyStrList, value, prefixWithDOT("")+key)
			}
		case reflect.Array, reflect.Slice:
			{
				objectValue := reflect.ValueOf(value)
				for index := 0; index < objectValue.Len(); index++ {
					propertyStrList = doMapToProperties(propertyStrList, objectValue.Index(index).Interface(), prefixWithDOT("")+key+"["+strconv.Itoa(index)+"]")
				}
			}
		case reflect.String:
			objectValue := reflect.ValueOf(value)
			objectValueStr := strings.ReplaceAll(objectValue.String(), "\n", "\\\n")
			propertyStrList = append(propertyStrList, prefixWithDOT("")+key+SignEqual+objectValueStr)
		default:
			propertyStrList = append(propertyStrList, prefixWithDOT("")+key+SignEqual+fmt.Sprintf("%v", value))
		}
	}
	resultStr := ""
	for _, propertyStr := range propertyStrList {
		resultStr += propertyStr + "\n"
	}

	return resultStr, nil
}

func doMapToProperties(propertyStrList []string, value interface{}, prefix string) []string {
	if nil == value {
		value = ""
	}

	valueKind := reflect.TypeOf(value).Kind()
	switch valueKind {
	case reflect.Map:
		{
			// map结构
			if reflect.ValueOf(value).Len() == 0 {
				return []string{}
			}

			for mapR := reflect.ValueOf(value).MapRange(); mapR.Next(); {
				mapKey := mapR.Key().Interface()
				mapValue := mapR.Value().Interface()
				propertyStrList = doMapToProperties(propertyStrList, mapValue, prefixWithDOT(prefix)+fmt.Sprintf("%v", mapKey))
			}
		}
	case reflect.Array, reflect.Slice:
		{
			objectValue := reflect.ValueOf(value)
			for index := 0; index < objectValue.Len(); index++ {
				propertyStrList = doMapToProperties(propertyStrList, objectValue.Index(index).Interface(), prefix+"["+strconv.Itoa(index)+"]")
			}
		}
	case reflect.String:
		objectValue := reflect.ValueOf(value)
		objectValueStr := strings.ReplaceAll(objectValue.String(), "\n", "\\\n")
		propertyStrList = append(propertyStrList, prefix+SignEqual+objectValueStr)
	default:
		objectValue := fmt.Sprintf("%v", reflect.ValueOf(value))
		propertyStrList = append(propertyStrList, prefix+SignEqual+objectValue)
	}
	return propertyStrList
}

func (p *yamlService) peelArray(nodeName string) (string, int) {
	var index = -1
	var name = nodeName
	var err error

	subData := rangePattern.FindAllStringSubmatch(nodeName, -1)
	if len(subData) > 0 {
		name = subData[0][1]
		indexStr := subData[0][2]
		if "" != indexStr {
			index, err = strconv.Atoi(indexStr)
			if err != nil {
				runtime.LogErrorf(p.ctx, "解析错误, nodeName=%s", nodeName)
				return "", -1
			}
		}
	}
	return name, index
}

func prefixWithDOT(prefix string) string {
	if "" == prefix {
		return ""
	}
	return prefix + "."
}

func (p *yamlService) wordToNode(lineWordList []string, nodeList []types.YamlNode, parentNode *types.YamlNode, lastNodeArrayFlag bool, index int, value string) ([]string, []types.YamlNode) {
	if len(lineWordList) == 0 {
		if lastNodeArrayFlag {
			node := types.YamlNode{Value: value, LastNodeIndex: -1}
			nodeList = append(nodeList, node)
		}
	} else {
		nodeName := lineWordList[0]
		nodeName, nextIndex := p.peelArray(nodeName)

		var node types.YamlNode
		if nil != parentNode {
			node = types.YamlNode{Name: nodeName, Parent: parentNode, LastNodeIndex: index}
		} else {
			node = types.YamlNode{Name: nodeName, LastNodeIndex: index}
		}
		lineWordList = lineWordList[1:]

		//如果节点下面的子节点数量为0，则为终端节点，也就是赋值节点
		if len(lineWordList) == 0 {
			if -1 == nextIndex {
				node.Value = value
			}
		}

		// nextIndex 不空，表示当前节点为数组，则之后的数据为他的节点数据
		if nextIndex != -1 {
			node.ArrayFlag = true
			var hasEqualsName = false

			//遍历查询节点是否存在
			for innerIndex := range nodeList {
				//如果节点名称已存在，则递归添加剩下的数据节点
				if nodeName == nodeList[innerIndex].Name && nodeList[innerIndex].ArrayFlag {
					yamlNodeIndex := nodeList[innerIndex].LastNodeIndex
					if -1 == yamlNodeIndex || index == yamlNodeIndex {
						hasEqualsName = true
						lineWordList, nodeList[innerIndex].ValueList = p.wordToNode(lineWordList, nodeList[innerIndex].ValueList, node.Parent, true, nextIndex, appendSpaceForArrayValue(value))
					}
				}
			}

			//如果遍历结果为节点名称不存在，则递归添加剩下的数据节点，并把新节点添加到上级yamlTree的子节点中
			if !hasEqualsName {
				lineWordList, node.ValueList = p.wordToNode(lineWordList, node.ValueList, node.Parent, true, nextIndex, appendSpaceForArrayValue(value))
				nodeList = append(nodeList, node)
			}
		} else {
			var hasEqualsName = false
			for innerIndex := range nodeList {
				if !lastNodeArrayFlag {
					//如果节点名称已存在，则递归添加剩下的数据节点
					if nodeName == nodeList[innerIndex].Name {
						hasEqualsName = true
						lineWordList, nodeList[innerIndex].Children = p.wordToNode(lineWordList, nodeList[innerIndex].Children, &nodeList[innerIndex], false, nextIndex, appendSpaceForArrayValue(value))
					}
				} else {
					//如果节点名称已存在，则递归添加剩下的数据节点
					if nodeName == nodeList[innerIndex].Name {
						yamlNodeIndex := nodeList[innerIndex].LastNodeIndex
						if -1 == yamlNodeIndex || index == yamlNodeIndex {
							hasEqualsName = true
							lineWordList, nodeList[innerIndex].Children = p.wordToNode(lineWordList, nodeList[innerIndex].Children, &nodeList[innerIndex], true, nextIndex, appendSpaceForArrayValue(value))
						}
					}
				}
			}

			//如果遍历结果为节点名称不存在，则递归添加剩下的数据节点，并把新节点添加到上级yamlTree的子节点中
			if !hasEqualsName {
				lineWordList, node.Children = p.wordToNode(lineWordList, node.Children, &node, false, nextIndex, appendSpaceForArrayValue(value))
				nodeList = append(nodeList, node)
			}
		}
	}
	return lineWordList, nodeList
}

func appendSpaceForArrayValue(value string) string {
	if !strings.HasPrefix(value, YamlNewLineDom) {
		return value
	}

	value = value[len(YamlNewLineDom):]
	valueTems := strings.Split(value, "\\n")

	strs := []string{}
	for _, element := range valueTems {
		tem := element
		if strings.HasSuffix(element, "\\") {
			tem = element[:len(element)-1]
		}
		strs = append(strs, IndentBlanks+tem)
	}
	return YamlNewLineDom + strings.Join(strs, "\n")
}

func GetPropertiesItemLineList(content string) []string {
	if "" == content {
		return []string{}
	}

	lineList := strings.Split(content, NewLine)
	var itemLineList []string
	var stringAppender string
	for _, line := range lineList {
		if strings.HasSuffix(content, "\\") {
			stringAppender += line + "\n"
		} else {
			stringAppender += line
			itemLineList = append(itemLineList, stringAppender)
			stringAppender = ""
		}
	}
	return itemLineList
}

func stringValueWrap(value string) string {
	if "" == value {
		return ""
	}
	// 对数组的数据进行特殊处理
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		return "'" + value + "'"
	}
	return value
}
