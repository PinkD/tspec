package tspec_test

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"testing"

	"github.com/stretchr/testify/suite"

	"tspec/samples"
	"tspec/tspec"
)

func TestTSpec(t *testing.T) {
	suite.Run(t, new(TSpecTestSuite))
}

type TSpecTestSuite struct {
	suite.Suite
	parser *tspec.Parser
	pkg    *ast.Package
}

func (s *TSpecTestSuite) SetupTest() {
	var err error
	s.parser = tspec.NewParser()
	s.pkg, err = s.parser.Import("tspec/samples")
	s.Require().NoError(err)
}

func (s *TSpecTestSuite) testParse(typeStr, assert string) {
	require := s.Require()

	schema, err := s.parser.Parse(s.pkg, typeStr)
	require.NoError(err)
	require.NotNil(schema)

	defs := s.parser.Definitions()
	bts, err := json.MarshalIndent(defs, "", "\t")
	require.NoError(err)
	require.Equal(string(bytes.TrimSpace(samples.MustAsset(assert))),
		string(bytes.TrimSpace(bts)))
	s.parser.Reset()
}

func (s *TSpecTestSuite) TestParse() {
	s.testParse("BasicTypes", "source/basic_types.json")
	s.testParse("NormalStruct", "source/normal_struct.json")
	s.testParse("StructWithNoExportField", "source/struct_with_no_export_field.json")
	s.testParse("StructWithAnonymousField", "source/struct_with_anonymous_field.json")
	s.testParse("StructWithCircularReference", "source/struct_with_circular_reference.json")
	s.testParse("StructWithInheritance", "source/struct_with_inheritance.json")
	s.testParse("MapType", "source/map_type.json")
	s.testParse("ArrayType", "source/array_type.json")
}

func (s *TSpecTestSuite) TestParseInvalidMap() {
	require := s.Require()

	schema, err := s.parser.Parse(s.pkg, "InvalidMap")
	require.Error(err)
	require.Nil(schema)
}
