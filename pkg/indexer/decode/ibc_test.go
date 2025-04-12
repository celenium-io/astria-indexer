// SPDX-FileCopyrightText: 2025 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package decode

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/celenium-io/astria-indexer/internal/storage"
	"github.com/stretchr/testify/require"
)

func Test_parseIbcMessages(t *testing.T) {
	tests := []struct {
		name string
		typ  string
		data string
	}{
		{
			name: "CreateClient",
			typ:  "/ibc.core.client.v1.MsgCreateClient",
			data: "CrEBCisvaWJjLmxpZ2h0Y2xpZW50cy50ZW5kZXJtaW50LnYxLkNsaWVudFN0YXRlEoEBCgluZXV0cm9uLTESBAgCEAMaBAiA3jQiBAiAvGkqAggoMgA6BwgBENSJ6whCGQoJCAEYASABKgEAEgwKAgABECEYBCAMMAFCGQoJCAEYASABKgEAEgwKAgABECAYASABMAFKB3VwZ3JhZGVKEHVwZ3JhZGVkSUJDU3RhdGVQAVgBEoUBCi4vaWJjLmxpZ2h0Y2xpZW50cy50ZW5kZXJtaW50LnYxLkNvbnNlbnN1c1N0YXRlElMKCwjt1OC7BhCaj8BLEiIKIJ+KuIJhVwXhxbJg9J+o6+R5zdbmkeSDS5xd+29sNFFQGiAFHQH1GteD52Qb2xN/9jLAsf7bKxeKWy8HMxB5YUn/tBotYXN0cmlhMTBud2VoajJrMDdweXh6bDc5NHBlNHU2YWx5dzhjZ2huY2t0ZTA1",
		}, {
			typ:  "/ibc.core.client.v1.MsgUpdateClient",
			name: "UpdateClient",
			data: "Cg8wNy10ZW5kZXJtaW50LTISpAcKJi9pYmMubGlnaHRjbGllbnRzLnRlbmRlcm1pbnQudjEuSGVhZGVyEvkGCtkECpoDCgIICxIRc3RyaWRlLWludGVybmFsLTEYgcX9BCIMCOeJ6LwGEJvwqcoCKkgKIHRJ//yecrMb62g2kDZpy84ubGGTwmwNnM6ReUBEVuOQEiQIARIgAUTogtYW+FxQjCfGCxAfS65QYU4tQ3otU32bmMF9deEyIGNsydA9BbrIcM345GH0Zk2SrYaILqXn3jf2arA3oMJlOiDjsMRCmPwcFJr79MiZb7kkJ65B5GSbk0yklZkbeFK4VUIgWnPjek3AHK1MwLFhloJrVItvR20eYr1OEurMuooi0pRKIFpz43pNwBytTMCxYZaCa1SLb0dtHmK9ThLqzLqKItKUUiAEgJG8fdwoP3e/v5HXPETaWMPfipy8hnQF2Lfz2q2iL1ogN0sUDiBWk26dcdezYXsIX48anhzVUzISG3LnJcfLupBiIOOwxEKY/BwUmvv0yJlvuSQnrkHkZJuTTKSVmRt4UrhVaiDjsMRCmPwcFJr79MiZb7kkJ65B5GSbk0yklZkbeFK4VXIUJG+P4/Otoo2tQBjjiv0s0ozL3ZUSuQEIgcX9BBpICiCGxDamzrdTvZTpEoi++OGZrFfvRoQ/ILDIeI1/t1d/LBIkCAESIFZGprCjGoEOp8dqhY2fj06y9m9XwT+WmMIYoaqtuCsDImgIAhIUJG+P4/Otoo2tQBjjiv0s0ozL3ZUaDAjsiei8BhD66bfRAiJA5Hy27KU00oYtqbsUj2lOlhj+cGxAVD7nGbphhl/reGANs48qKhwALURznebHcMBvNQZqHJLG3qa3wKp+J/qGDxKHAQo/ChQkb4/j862ija1AGOOK/SzSjMvdlRIiCiDyz35EQUnFVkYo8ltJ+K3K56cX4a4J9Qq2ypH43yrUWBjAjbcBEj8KFCRvj+PzraKNrUAY44r9LNKMy92VEiIKIPLPfkRBScVWRijyW0n4rcrnpxfhrgn1CrbKkfjfKtRYGMCNtwEYwI23ARoHCAEQ3Kr5BCKHAQo/ChQkb4/j862ija1AGOOK/SzSjMvdlRIiCiDyz35EQUnFVkYo8ltJ+K3K56cX4a4J9Qq2ypH43yrUWBjAjbcBEj8KFCRvj+PzraKNrUAY44r9LNKMy92VEiIKIPLPfkRBScVWRijyW0n4rcrnpxfhrgn1CrbKkfjfKtRYGMCNtwEYwI23ARotYXN0cmlhMWt6ZXFqYTlnazgzZHo3cnI3ZjNwNXJ6dGpzYTUweXluaHJla2Z2",
		}, {
			name: "UpgradeClient",
			typ:  "/ibc.core.client.v1.MsgUpgradeClient",
			data: "",
		}, {
			// 	name: "SubmitMisbehaviour",
			// 	typ:  "/ibc.core.client.v1.MsgSubmitMisbehaviour",
			// 	data: "",
			// }, {
			name: "ConnectionOpenInit",
			typ:  "/ibc.core.connection.v1.MsgConnectionOpenInit",
			data: "ChAwNy10ZW5kZXJtaW50LTIyEhoKETA3LXRlbmRlcm1pbnQtMTUwGgUKA2liYxojCgExEg1PUkRFUl9PUkRFUkVEEg9PUkRFUl9VTk9SREVSRUQqLWFzdHJpYTEwbndlaGoyazA3cHl4emw3OTRwZTR1NmFseXc4Y2dobmNrdGUwNQ==",
		}, {
			name: "ConnectionOpenTry",
			typ:  "/ibc.core.connection.v1.MsgConnectionOpenTry",
			data: "",
		}, {
			name: "ConnectionOpenAck",
			typ:  "/ibc.core.connection.v1.MsgConnectionOpenAck",
			data: "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",
		}, {
			name: "ConnectionOpenConfirm",
			typ:  "/ibc.core.connection.v1.MsgConnectionOpenConfirm",
			data: "",
		}, {
			name: "ChannelOpenInit",
			typ:  "/ibc.core.channel.v1.MsgChannelOpenInit",
			data: "Cgh0cmFuc2ZlchInCAEQARoKCgh0cmFuc2ZlciIMY29ubmVjdGlvbi00KgdpY3MyMC0xGi1hc3RyaWExMG53ZWhqMmswN3B5eHpsNzk0cGU0dTZhbHl3OGNnaG5ja3RlMDU=",
		}, {
			name: "ChannelOpenTry",
			typ:  "/ibc.core.channel.v1.MsgChannelOpenTry",
			data: "",
		}, {
			name: "ChannelOpenAck",
			typ:  "/ibc.core.channel.v1.MsgChannelOpenAck",
			data: "Cgh0cmFuc2ZlchIJY2hhbm5lbC00GgxjaGFubmVsLTYyMzYiB2ljczIwLTEqpQwK+AkK9QkKMGNoYW5uZWxFbmRzL3BvcnRzL3RyYW5zZmVyL2NoYW5uZWxzL2NoYW5uZWwtNjIzNhI0CAIQARoVCgh0cmFuc2ZlchIJY2hhbm5lbC00Ig5jb25uZWN0aW9uLTExMCoHaWNzMjAtMRoOCAEYASABKgYAAoaU1hEiLAgBEigCBIaU1hEgPTNybQhcAkrPMUdkNm3YwqwqyrDGa3g1PNp+BxQNW9cgIiwIARIoBAiGlNYRIFd79qnwRKc6R/o4XXuWh+W3I0vMF9aCUzUcrewhB6YwICIuCAESBwYQhpTWESAaISBDBPGCdBJd1f4Ubtu6Lv2CPRA2ri8TlTpGdXDwUnhmOSIsCAESKAgehpTWESCayGP3xZVIOOTDpqU9+iZ1ITszK+QNDNpiVNEDY1MoJSAiLggBEgcKNIaU1hEgGiEgt/7IGmRgrRTH8/ncvG8PiGGQiVI+HqEcATlqvbBxqpAiLQgBEikOlAGGlNYRIHOsWcWvPFsl3ddTYCjCfznfiQdsHVvY/0hyeaI6ZdUKICIvCAESCBD4AYaU1hEgGiEgC5UyhU5f2gTSjO1I8mz7wBGXcoD99n6ggJp56/duNSYiLwgBEggSzgOGlNYRIBohIH9nBt++riL1umFQm6/HmkL8IxduaJengtwkJQ6UWj8iIi0IARIpFLIHhpTWESBe6GAOCmXtHa0oae4sa8riiYJPilbTLhkgrWMuki4xtCAiLQgBEikYkBOGlNYRIKKbbfxJ7s4PuYrTmmpProBq5/TYB4veWAIUqaaO2VMNICItCAESKRrqIoaU1hEgArPmB0izJUr5lL8kXFATG3iyNTsggT9ZVoMjl6IQLTMgIi0IARIpHIA4hpTWESCZvjd798TkrZhiABlLOLutIJcEfZcYOKzt36gtdWSbwSAiLQgBEike2FyGlNYRICnAP8sp9knwIfq2usqlsC8GE4lXDXEdeeWbpGB+MizhICIuCAESKiDmrAGGlNYRICGWI6NH0Wc6AL2TX9gNEZFz4089XYdi/J68Vs74Sln+ICIuCAESKiKAtgKGlNYRIPeNRSFGBM2ZDThyTV/n8HwmLunkHjjGKSDM5JHhXTAaICIuCAESKiTa9QWGlNYRIATMOdT1ONw9mPQ6r3czkvqvxmBeNaS0qni0nV1/an9dICIuCAESKibqkAqGlNYRIDZhkY0jQFV73KuA22P5Ujo/Kgb/zSrxqol1BocDD82GICIuCAESKij26RKGlNYRIM05BmziTntylIn9BfqF87vZDYZ0kwLT2nAg/6x68qW9ICIuCAESKiqMgh+GlNYRINTZVMIR+/YGEvMwMbDsPo23UpX2UsKXfvubp/9fsIo6ICIuCAESKiycljCGlNYRIFuLD3yM8U37c+nRs0b1Gn8fX5Jxym26smsJ2YYR3G1uICIuCAESKi6wxnmGlNYRICeaIzikVdvL+EcOkEZmtqkAHQZr0j0LfOhRUz6E9kJgICIxCAESCjD81ucBipTWESAaISAJvmmtGJurf12oOJEYoYcRLmuUHWbjPacjuHwHmlLzzSIxCAESCjKyqbIDipTWESAaISBtztiDhEKAVO2EJtmXOC49uASk4mjstCXORzlpraEg6iIvCAESKzT606UGipTWESBAZ5xPI/QFycEkGJVUHCJLXHwxZqRrrcoRLRzIgvJmfiAKpwIKpAIKA2liYxIgEnu7vs3vfYmy1HlcINv37gOcRXMk7AzIsA690KFQ1ZUaCQgBGAEgASoBACInCAESAQEaIG7a07baeY1MZctEUFL7MeW87Kzu6tualuZ5LL3ceSfxIicIARIBARogxZoclubk6SoWRBzBpY6f482gZA/Xa0m71jlVXjiiNc8iJQgBEiEBjsyu4n00iM5xq/h/JuMZ5A1QYMxx2oesLvPsKN799I8iJwgBEgEBGiDd0D0udrMe6GlV5qNH2Suqok6bqmBzQ1KjhVj3FZkOECIlCAESIQEUntLmKbcQKJ6slw2j4a4re1jNkopXk9iMUmfImTOm2yInCAESAQEaIH9awFy1j8sl1bb8042uW8c++jPO8h6AM3PrtdckswacMgcIARCGiusIOi1hc3RyaWExMG53ZWhqMmswN3B5eHpsNzk0cGU0dTZhbHl3OGNnaG5ja3RlMDU=",
		}, {
			name: "ChannelOpenConfirm",
			typ:  "/ibc.core.channel.v1.MsgChannelOpenConfirm",
			data: "",
		}, {
			name: "ChannelCloseInit",
			typ:  "/ibc.core.channel.v1.MsgChannelCloseInit",
			data: "",
		}, {
			name: "ChannelCloseConfirm",
			typ:  "/ibc.core.channel.v1.MsgChannelCloseConfirm",
			data: "",
		}, {
			name: "RecvPacket",
			typ:  "/ibc.core.channel.v1.MsgRecvPacket",
			data: "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",
		}, {
			name: "Acknowledgement",
			typ:  "/ibc.core.channel.v1.MsgAcknowledgement",
			data: "CrEECMk3Egh0cmFuc2ZlchoJY2hhbm5lbC0wIgh0cmFuc2ZlcioKY2hhbm5lbC00ODLeA3siZGVub20iOiJ0cmFuc2Zlci9jaGFubmVsLTAvdXRpYSIsImFtb3VudCI6IjIwMDAwMDAiLCJzZW5kZXIiOiJhc3RyaWExM3ZwdGRhZnl0dHBtbHdwcHQwczg0NGVmZXkyY3BjMG1ldnk5MnAiLCJyZWNlaXZlciI6ImNlbGVzdGlhMTZhd3ZhcGQwNWRuenFmZWFrc3g2MjZ6dGN1eDh3Z3pmNWw1YWdoIiwibWVtbyI6IntcInJvbGx1cEJsb2NrTnVtYmVyXCI6XCIzODE3NTg2XCIsXCJyb2xsdXBXaXRoZHJhd2FsRXZlbnRJZFwiOlwiMHg5MTkxNjcyZjI2YmM3MmMyOWVhODEwMjRjNGRiY2E0YThhNmI0NGUyZjFkZjU2NjBlMGUxNDZiMzUxMTk3NzA0LjB4MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMFwiLFwicm9sbHVwUmV0dXJuQWRkcmVzc1wiOlwiMHgwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDA1MzY3MTA5YjhhYjU3YTViNWZjNGQzMDIxYTc2OWJjMjY4OGNlNTczXCJ9In06Fgj///////////8BEP///////////wFAt46TivLJwJAYEhF7InJlc3VsdCI6IkFRPT0ifRqdDAqZCgqWCgo2YWNrcy9wb3J0cy90cmFuc2Zlci9jaGFubmVscy9jaGFubmVsLTQ4L3NlcXVlbmNlcy83MTEzEiAI91V+1Rgm/hjYRRK/JOx1AB7bryEjpHffcqCp82QKfBoOCAEYASABKgYAArCF2AMiLAgBEigCBLCF2AMgz09qw8MnoJ57eA6f+sVmiq5+xRuCH62mG0bsdiAXUEMgIiwIARIoBAawhdgDIESlydHykNSulQzUVHU6G614ULQRMWyLbKbNyWo1NglSICIsCAESKAYOsIXYAyB3o/KrdP87AtRrxVSsAyARufeH3bfV38qQMjnQzhtjoyAiLggBEgcIFLCF2AMgGiEgpPe6VjF5ETHdbgeXda182qa/Zch7VsMs6nl1UCKYAiAiLAgBEigKKrCF2AMgxFHDEc1DQwCxx/TKBwtM5IvWzdnoJa7C5mO4EdrKV50gIiwIARIoDEKwhdgDIF92YXejSl9R3rrzkxbonllycoru4+cS58Ib3b13pD+0ICItCAESKQ6iAbCF2AMgHRyCxLwadxC7mxPl2hfoV79K0jJLYjzAq5TR6h7eN0EgIi0IARIpENwCsIXYAyBvdyttqSoMfAIrKjECMrEiMIimqJ9RSTW7TkAboxuV1yAiLQgBEikSnASwhdgDIH0rf9I3koPxCI5dY3hMggUKck2SRUh2Ec7LrU7le7M8ICItCAESKRT0BrCF2AMgfDk5kfLAjE2barvURV/XTxiu9Cc9NkPMWYtUYVjEwpwgIi8IARIIFswLsIXYAyAaISBR/f73gzHVQNxWql3rql1t878mhl9GT3CmXNIEvKqo7CIvCAESCBiiGbCF2AMgGiEgsEHcbs/9ZfUhgTjzXzntggzraoJ+qbK3c4GFpmixZwIiLQgBEika9C6whdgDIM5j43MHBYcqAdS/Ip3aIziolrAywvRbNk9GRPr+Rc3mICItCAESKRz8RLCF2AMgdf9Hs4lQGEvB0NHDFrPElBqu+nEdVAEXtw6j1RqzqN8gIi0IARIpHqpvsIXYAyCML1Vnn/PR/Hce3ecIlTvTAI0JQ6hckQLCrTVNxDP62yAiLggBEiog4oICsIXYAyA+J3HyQsgtz72CcCpEGmPtbOqGPyNGiZl9pjDksHQMBiAiLggBEioi/vYDsIXYAyB3k9vi4yzkUIwEGwnjMowbeH6Xy2yXDAoSF4bmSGo5xSAiLggBEiokhN0HsIXYAyBKvxMfEZd6DVdg31wqFhw8KB2FqJNFVszpfWck2wBUAyAiMAgBEgkmuIwLsIXYAyAaISD1PDuKfJAKIk+2tLcfbDF4Ef9MkTz6eVKjdNx53+VJviIuCAESKiie/Q+whdgDIKy/VVTuOrogl3XtDzQvJdbxfTQGQVSUrDWeEncBzC8tICIwCAESCSrewCawhdgDIBohIK+oljffktmIS1GAvy7yw9llGqa0k+c/1+kV2QqfOPVCIi4IARIqLP7CSrCF2AMg31qli9bkPODAYKAJ4kwAdz3EinDuND0woBAqEPuMNu0gIjEIARIKMNCL3AGwhdgDIBohIMYWMZZpBgJ7btsL+86YGqlRxZclhOyMZ5SKawcJenN6Ii8IARIrMrSH4AKwhdgDIHePiW9eaAYWu3pEHYb6WliJIKQMrIqmRKyq80JUhOeRICIxCAESCjScjOUFsIXYAyAaISD+sf2btfAX1cij6CAFcd4ibHj63QKCaRqYOJMFPvRR+Ar+AQr7AQoDaWJjEiCNa0IqRGauEhrdQap/6lAU1EkdqNxRwhS2vGm97FZdABoJCAEYASABKgEAIiUIARIhAf8dTzAYsEMU561Nr6qMM2jDrmF8+y4H5IdtEKMSkqeEIicIARIBARogeQb/KDAMcBY+Vsw9Gdlt8WhryKOy8VSBjRpM2NnuAKgiJwgBEgEBGiA7E6EKzWflHY5oIUyU0JKveREbQHtC7glG90YYe1fRKSIlCAESIQEo/R10iZrm9DoDWfngr0SB/+Ffka2SqBrdgRtTWEb3gCInCAESAQEaIN/YbJEAHOdDrcS6gECqPwejnTz3Su7Sa0ll5J7aE5MMIgUQ2YLsASotYXN0cmlhMW5ybnd0cjVlbXNnOTQwdmx0aGg2eXZnNzBsbmNrZGdwejBzdmhh",
		}, {
			name: "Timeout",
			typ:  "/ibc.core.channel.v1.MsgTimeout",
			data: "CrIECKc3Egh0cmFuc2ZlchoJY2hhbm5lbC0wIgh0cmFuc2ZlcioKY2hhbm5lbC00ODLfA3siZGVub20iOiJ0cmFuc2Zlci9jaGFubmVsLTAvdXRpYSIsImFtb3VudCI6IjMwMDE0MDAwIiwic2VuZGVyIjoiYXN0cmlhMTN2cHRkYWZ5dHRwbWx3cHB0MHM4NDRlZmV5MmNwYzBtZXZ5OTJwIiwicmVjZWl2ZXIiOiJjZWxlc3RpYTEwYW5oeHFnbGVqbXg3YXBjNW40amswbmx0Z3JkcmY0ZDdtaDhwMCIsIm1lbW8iOiJ7XCJyb2xsdXBCbG9ja051bWJlclwiOlwiMzgwNTU1N1wiLFwicm9sbHVwV2l0aGRyYXdhbEV2ZW50SWRcIjpcIjB4YTgyYzZkNmE1M2EwNGFiMTk5MzVlN2ExMDFiZjIxMzNkN2QzMGMzOTFmNGVmOTFlOTQyNzA4NTY1ODlmOWVkYy4weDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDBcIixcInJvbGx1cFJldHVybkFkZHJlc3NcIjpcIjB4MDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwNTNiYTZkY2JkYzc5MTVjYTk4NTM2NzIyZDc3ZDY0MzkxMTU4N2JiYlwifSJ9OhYI////////////ARD///////////8BQOL4vc3ZqrqQGBKVFgqRFBKOFAo6cmVjZWlwdHMvcG9ydHMvdHJhbnNmZXIvY2hhbm5lbHMvY2hhbm5lbC00OC9zZXF1ZW5jZXMvNzA3ORL9CQo6cmVjZWlwdHMvcG9ydHMvdHJhbnNmZXIvY2hhbm5lbHMvY2hhbm5lbC00OC9zZXF1ZW5jZXMvNzA3OBIBARoOCAEYASABKgYAArK21wMiLAgBEigCBLK21wMgxKzgml+m1kSxcWZRAAtdi9LzXHOKugZMkHq1kuAhqCAgIiwIARIoBAayttcDIIGn3gWkF9SSewHhRF3qUt2IL05u3d34HszBh4GONnATICIsCAESKAYOsrbXAyALuwc5q+V9rZnhsMYz4ydKq7ZOLrszkw6lD1i379qNQyAiLggBEgcIFrK21wMgGiEgvRggj3OyK+J2w40GTm8Vd/pfkvTLAXmc/AlMIOkUs6MiLAgBEigKIrK21wMgcVvSjahz86briojzFV9j9/GPjEsAlqn7GdqF6BmX5hAgIiwIARIoDFKyttcDIGxKqnyJgIMtFA1ZlooOfl2p1xgDlC7QJzDfC0zaTocVICItCAESKQ6mAbK21wMgcKOKz1gQP7jkz9eKg6StyRxmN3IYkyRarESOJ3t5Y48gIi8IARIIEMoCsrbXAyAaISAGSzXoeQ+38C8tCPCd+3W2C8//4JCJb0Htm0KRHV6WjSIvCAESCBLkA7K21wMgGiEg0IM1CCGA9BJ36f5qX6vxiX1IicdgWgsz/pyIoqVNO+UiLwgBEggU2geyttcDIBohIBMOtVRsEsK0g31/Wz3Y+X8sjN4PjkH9YUmni3NRHV6jIi0IARIpFtoSsrbXAyDWXwFw4F5dnH7NVzyFIRl2B7o05ZIewJdR3Ta2K5QhCyAiLQgBEikYqh6yttcDIDHF4dDdoYo9QHnCKfPki9IlM56OQ5xTMdwGyMuEwM+dICIvCAESCBqcKbK21wMgGiEgB8ve8itdZtXjxrjKupxh//u4rgZGmoC7XWD5zruD5/ciLQgBEikcolOyttcDIMGNByWM9ArEUlTh6chuzQ3w0eQpjcBgce+3QNRqKkifICItCAESKR68frK21wMgPI/xGL1YxKesw063IqmZWEvXvG1qTiHTaUNBEEgZT6kgIjAIARIJIJbXAbK21wMgGiEgOOx78HVaF/0WyRG+qcS8HOM5pBNr2S6xTXCzfS/0pP8iLggBEioiytQDsrbXAyAuDSA68B3HsWEOnE/19OJZUg5gkUlxSj4ix8sfTNJcryAiLggBEiokpukHsrbXAyB6KeKEAVNMNqaPS1ICGKs9mLIys5kTxA3fa9oR8mM7XiAiLggBEiom0qALsrbXAyDe+y1+QxTjN8R6Sd4KtCwRwqfRUlHPjV4mmfZ3mtG0iiAiLggBEioohrQUiLvXAyA7EA7C5J319g7ao2o97SoWfZF7dmUC0Pqb41hkcK5bJyAiLggBEioq9pwtiLvXAyAGvmgNG7tU2b6tb8EfWLQYbaxxL0rlvPoS4Etqo4UjTCAiMAgBEgksivdMiLvXAyAaISDVuyBEsjBSUum4Twru4s3UVYfOYC352zIRq84NV7gCoSIxCAESCjDav6wBiLvXAyAaISDrbyEa9VXwGezoYyrv7YvGuF4YWRIXDDaY0pD83AJN2CIvCAESKzK4+IQDjrvXAyCQRKxwW+t0J/j91qoo0ax083gixJk5sgSuId7DPzcahyAiLwgBEis0jPbkBY671wMgQ8R3l4GpAFFmX5fdUzFXSPHATqk0FCyKuaT2/fkt+ZwgGs8JCjhyZWNlaXB0cy9wb3J0cy90cmFuc2Zlci9jaGFubmVscy9jaGFubmVsLTQ4L3NlcXVlbmNlcy83MRIBARoOCAEYASABKgYAApjgxwIiLggBEgcCBKThxwIgGiEgExZysqPYdANpwKzJRGtsKuWtZJ2O9ZKYOyRT0Qb9FXIiLggBEgcECMLz1gMgGiEg49tvv4ZXCN1ajkKU+C6HXn9g5z44fW6wG7Inie7N1XQiLAgBEigIFrK21wMgWn+Bq12HhkzFSHtpJBEdLn1nvGi4csaJ5XbFgYM/TfggIiwIARIoCiKyttcDIHFb0o2oc/Om64qI8xVfY/fxj4xLAJap+xnahegZl+YQICIsCAESKAxSsrbXAyBsSqp8iYCDLRQNWZaKDn5dqdcYA5Qu0Ccw3wtM2k6HFSAiLQgBEikOpgGyttcDIHCjis9YED+45M/XioOkrckcZjdyGJMkWqxEjid7eWOPICIvCAESCBDKArK21wMgGiEgBks16HkPt/AvLQjwnft1tgvP/+CQiW9B7ZtCkR1elo0iLwgBEggS5AOyttcDIBohINCDNQghgPQSd+n+al+r8Yl9SInHYFoLM/6ciKKlTTvlIi8IARIIFNoHsrbXAyAaISATDrVUbBLCtIN9f1s92Pl/LIzeD45B/WFJp4tzUR1eoyItCAESKRbaErK21wMg1l8BcOBeXZx+zVc8hSEZdge6NOWSHsCXUd02tiuUIQsgIi0IARIpGKoesrbXAyAxxeHQ3aGKPUB5winz5IvSJTOejkOcUzHcBsjLhMDPnSAiLwgBEgganCmyttcDIBohIAfL3vIrXWbV48a4yrqcYf/7uK4GRpqAu11g+c67g+f3Ii0IARIpHKJTsrbXAyDBjQcljPQKxFJU4enIbs0N8NHkKY3AYHHvt0DUaipInyAiLQgBEikevH6yttcDIDyP8Ri9WMSnrMNOtyKpmVhL17xtak4h02lDQRBIGU+pICIwCAESCSCW1wGyttcDIBohIDjse/B1Whf9FskRvqnEvBzjOaQTa9kusU1ws30v9KT/Ii4IARIqIsrUA7K21wMgLg0gOvAdx7FhDpxP9fTiWVIOYJFJcUo+IsfLH0zSXK8gIi4IARIqJKbpB7K21wMgeinihAFTTDamj0tSAhirPZiyMrOZE8QN32vaEfJjO14gIi4IARIqJtKgC7K21wMg3vstfkMU4zfEekneCrQsEcKn0VJRz41eJpn2d5rRtIogIi4IARIqKIa0FIi71wMgOxAOwuSd9fYO2qNqPe0qFn2Re3ZlAtD6m+NYZHCuWycgIi4IARIqKvacLYi71wMgBr5oDRu7VNm+rW/BH1i0GG2scS9K5bz6EuBLaqOFI0wgIjAIARIJLIr3TIi71wMgGiEg1bsgRLIwUlLpuE8K7uLN1FWHzmAt+dsyEavODVe4AqEiMQgBEgow2r+sAYi71wMgGiEg628hGvVV8Bns6GMq7+2LxrheGFkSFww2mNKQ/NwCTdgiLwgBEisyuPiEA4671wMgkESscFvrdCf4/daqKNGsdPN4IsSZObIEriHewz83GocgIi8IARIrNIz25AWOu9cDIEPEd5eBqQBRZl+X3VMxV0jxwE6pNBQsirmk9v35LfmcIAr+AQr7AQoDaWJjEiCng/ZNDASxBHD0sGpIeaKq7GITlJRi/VBqHW/ysCNVPRoJCAEYASABKgEAIiUIARIhAf8dTzAYsEMU561Nr6qMM2jDrmF8+y4H5IdtEKMSkqeEIicIARIBARogeQb/KDAMcBY+Vsw9Gdlt8WhryKOy8VSBjRpM2NnuAKgiJwgBEgEBGiDasiSa6M9g9YqCAtZx+zEC5nhNfJ1xKjdfFp9s6TSjjCIlCAESIQFLlifdf1WZM5Wwa87RntzeKLe7is7fOUxGMv9xcz++RyInCAESAQEaIBcIfEwvOGu9km32pZLmaewF9qNkA9oLi+I5ro5FP9xNGgUQyN3rASCnNyotYXN0cmlhMW5ybnd0cjVlbXNnOTQwdmx0aGg2eXZnNzBsbmNrZGdwejBzdmhh",
		}, {
			name: "Unknown",
			typ:  "",
			data: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.data == "" {
				return
			}
			action := storage.Action{
				Data: make(map[string]any),
			}
			ctx := NewContext(nil, time.Now())

			b, err := base64.StdEncoding.DecodeString(tt.data)
			require.NoError(t, err)

			err = parseIbcMessages(tt.typ, b, &action, &ctx)
			require.NoError(t, err)
			require.Contains(t, action.Data, "msg")
			require.NotEmpty(t, action.Data["msg"])
		})
	}
}
