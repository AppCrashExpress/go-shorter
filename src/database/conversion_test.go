package database

import "testing"

func TestConversion(t *testing.T) {
    type Testcase struct {
        input    string;
        output   string;
        descript string;
    }

    testcases := []Testcase {
        {
            "",
            "2jmj7l5rSw",
            "Test empty input",
        },
        {
            "some string",
            "i0XkvRxqy4",
            "Test some input",
        },
        {
            "fad",
            "2Dn5FVwM0c",
            "Test with dash",
        },
    }

    for _, test := range testcases {
        t.Run(test.descript, func(t *testing.T) {
            result := string(ConvertUrl([]byte(test.input)))
            if !(test.output == result) {
                t.Errorf("Failed: %v\nExpected: %v\nGot: %v", test.descript, test.output, result)
            }
        })
    }

}


