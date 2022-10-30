package service

import (
	"context"
	"encoding/json"
	"github.com/aaronchen2k/deeptest/internal/pkg/helper/openapi"
	_fileUtils "github.com/aaronchen2k/deeptest/pkg/lib/file"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
)

type SpecService struct {
}

func (s *SpecService) Load(path, typ string) (content string, err error) {
	var doc3 *openapi3.T

	if typ == "openapi2" {
		v2Content := _fileUtils.ReadFileBuf(path)

		var doc2 openapi2.T
		err = json.Unmarshal(v2Content, &doc2)
		if err != nil {
			return
		}

		doc3, err = openapi.ToV3(&doc2)
		if err != nil {
			return
		}

		//err = doc3.Validate(context.Background())

		var bytes []byte
		bytes, err = json.MarshalIndent(doc3, "", "\t")
		if err != nil {
			return
		}

		content = string(bytes)
	} else if typ == "openapi3" {
		ctx := context.Background()
		loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}

		doc3, err = loader.LoadFromFile(path)
		if err != nil {
			return
		}

		//_ := doc.Validate(ctx)

		var bytes []byte
		bytes, err = json.Marshal(doc3)
		if err != nil {
			return
		}

		content = string(bytes)
	}

	return
}
