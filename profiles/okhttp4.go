package profiles

import (
	"github.com/nukilabs/http/http2"
	tls "github.com/nukilabs/utls"
)

func Okhttp4Android(sdk int) ClientProfile {
	switch sdk {
	case 21:
		return Okhttp4Android21
	case 22:
		return Okhttp4Android22
	case 23:
		return Okhttp4Android23
	case 24:
		return Okhttp4Android24
	case 25:
		return Okhttp4Android25
	case 26:
		return Okhttp4Android26
	case 27:
		return Okhttp4Android27
	case 28:
		return Okhttp4Android28
	case 29:
		return Okhttp4Android29
	case 30:
		return Okhttp4Android30
	case 31:
		return Okhttp4Android31
	case 32:
		return Okhttp4Android32
	case 33:
		return Okhttp4Android33
	case 34:
		return Okhttp4Android34
	default:
		return Okhttp4Android29
	}
}

var Okhttp4Android21 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.FAKE_TLS_EMPTY_RENEGOTIATION_INFO_SCSV,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.SNIExtension{},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00, 0x01, 0x02,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.CurveID(14),
					tls.CurveID(13),
					tls.CurveP521,
					tls.CurveID(11),
					tls.CurveID(12),
					tls.CurveP384,
					tls.CurveID(9),
					tls.CurveID(10),
					tls.CurveID(22),
					tls.CurveP256,
					tls.CurveID(8),
					tls.CurveID(6),
					tls.CurveID(7),
					tls.CurveID(20),
					tls.CurveID(21),
					tls.CurveID(4),
					tls.CurveID(5),
					tls.CurveID(18),
					tls.CurveID(19),
					tls.CurveID(1),
					tls.CurveID(2),
					tls.CurveID(3),
					tls.CurveID(15),
					tls.CurveID(16),
					tls.CurveID(17),
				}},
				&tls.SessionTicketExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.PKCS1WithSHA512,
					0x602,
					tls.ECDSAWithP521AndSHA512,
					tls.PKCS1WithSHA384,
					0x502,
					tls.ECDSAWithP384AndSHA384,
					tls.PKCS1WithSHA256,
					0x402,
					tls.ECDSAWithP256AndSHA256,
					0x301,
					0x302,
					0x303,
					tls.PKCS1WithSHA1,
					0x202,
					tls.ECDSAWithSHA1,
				}},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android22 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.FAKE_TLS_EMPTY_RENEGOTIATION_INFO_SCSV,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.SNIExtension{},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00, 0x01, 0x02,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.CurveID(14),
					tls.CurveID(13),
					tls.CurveP521,
					tls.CurveID(11),
					tls.CurveID(12),
					tls.CurveP384,
					tls.CurveID(9),
					tls.CurveID(10),
					tls.CurveID(22),
					tls.CurveP256,
					tls.CurveID(8),
					tls.CurveID(6),
					tls.CurveID(7),
					tls.CurveID(20),
					tls.CurveID(21),
					tls.CurveID(4),
					tls.CurveID(5),
					tls.CurveID(18),
					tls.CurveID(19),
					tls.CurveID(1),
					tls.CurveID(2),
					tls.CurveID(3),
					tls.CurveID(15),
					tls.CurveID(16),
					tls.CurveID(17),
				}},
				&tls.SessionTicketExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.PKCS1WithSHA512,
					0x602,
					tls.ECDSAWithP521AndSHA512,
					tls.PKCS1WithSHA384,
					0x502,
					tls.ECDSAWithP384AndSHA384,
					tls.PKCS1WithSHA256,
					0x402,
					tls.ECDSAWithP256AndSHA256,
					0x301,
					0x302,
					0x303,
					tls.PKCS1WithSHA1,
					0x202,
					tls.ECDSAWithSHA1,
				}},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android23 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.FAKE_TLS_EMPTY_RENEGOTIATION_INFO_SCSV,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.SessionTicketExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.PKCS1WithSHA512,
					tls.ECDSAWithP521AndSHA512,
					tls.PKCS1WithSHA384,
					tls.ECDSAWithP384AndSHA384,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP256AndSHA256,
					0x301,
					0x303,
					tls.PKCS1WithSHA1,
					tls.ECDSAWithSHA1,
				}},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.CurveP256,
					tls.CurveP384,
					tls.CurveP521,
				}},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android24 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.SessionTicketExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.PKCS1WithSHA512,
					tls.ECDSAWithP521AndSHA512,
					tls.PKCS1WithSHA384,
					tls.ECDSAWithP384AndSHA384,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP256AndSHA256,
					0x301,
					0x303,
					tls.PKCS1WithSHA1,
					tls.ECDSAWithSHA1,
				}},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.CurveP256,
				}},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android25 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.SessionTicketExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.PKCS1WithSHA512,
					tls.ECDSAWithP521AndSHA512,
					tls.PKCS1WithSHA384,
					tls.ECDSAWithP384AndSHA384,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP256AndSHA256,
					0x301,
					0x303,
					tls.PKCS1WithSHA1,
					tls.ECDSAWithSHA1,
				}},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.CurveP256,
					tls.CurveP384,
					tls.CurveP521,
				}},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android26 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.SessionTicketExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PKCS1WithSHA384,
					tls.ECDSAWithP521AndSHA512,
					tls.PKCS1WithSHA512,
					tls.PKCS1WithSHA1,
				}},
				&tls.StatusRequestExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android27 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.SessionTicketExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PKCS1WithSHA384,
					tls.ECDSAWithP521AndSHA512,
					tls.PKCS1WithSHA512,
					tls.PKCS1WithSHA1,
				}},
				&tls.StatusRequestExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android28 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.SessionTicketExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PSSWithSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PSSWithSHA384,
					tls.PKCS1WithSHA384,
					tls.PSSWithSHA512,
					tls.PKCS1WithSHA512,
					tls.PKCS1WithSHA1,
				}},
				&tls.StatusRequestExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android29 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SessionTicketExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.StatusRequestExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PSSWithSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PSSWithSHA384,
					tls.PKCS1WithSHA384,
					tls.PSSWithSHA512,
					tls.PKCS1WithSHA512,
					tls.PKCS1WithSHA1,
				}},
				&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
					{Group: tls.X25519},
				}},
				&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
					tls.PskModeDHE,
				}},
				&tls.SupportedVersionsExtension{Versions: []uint16{
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android30 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SessionTicketExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.StatusRequestExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PSSWithSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PSSWithSHA384,
					tls.PKCS1WithSHA384,
					tls.PSSWithSHA512,
					tls.PKCS1WithSHA512,
					tls.PKCS1WithSHA1,
				}},
				&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
					{Group: tls.X25519},
				}},
				&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
					tls.PskModeDHE,
				}},
				&tls.SupportedVersionsExtension{Versions: []uint16{
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android31 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SessionTicketExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.StatusRequestExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PSSWithSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PSSWithSHA384,
					tls.PKCS1WithSHA384,
					tls.PSSWithSHA512,
					tls.PKCS1WithSHA512,
					tls.PKCS1WithSHA1,
				}},
				&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
					{Group: tls.X25519},
				}},
				&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
					tls.PskModeDHE,
				}},
				&tls.SupportedVersionsExtension{Versions: []uint16{
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android32 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SessionTicketExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.StatusRequestExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PSSWithSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PSSWithSHA384,
					tls.PKCS1WithSHA384,
					tls.PSSWithSHA512,
					tls.PKCS1WithSHA512,
					tls.PKCS1WithSHA1,
				}},
				&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
					{Group: tls.X25519},
				}},
				&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
					tls.PskModeDHE,
				}},
				&tls.SupportedVersionsExtension{Versions: []uint16{
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android33 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SessionTicketExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.StatusRequestExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PSSWithSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PSSWithSHA384,
					tls.PKCS1WithSHA384,
					tls.PSSWithSHA512,
					tls.PKCS1WithSHA512,
					tls.PKCS1WithSHA1,
				}},
				&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
					{Group: tls.X25519},
				}},
				&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
					tls.PskModeDHE,
				}},
				&tls.SupportedVersionsExtension{Versions: []uint16{
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}

var Okhttp4Android34 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				0x00,
			},
			Extensions: []tls.TLSExtension{
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SessionTicketExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.StatusRequestExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PSSWithSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PSSWithSHA384,
					tls.PKCS1WithSHA384,
					tls.PSSWithSHA512,
					tls.PKCS1WithSHA512,
					tls.PKCS1WithSHA1,
				}},
				&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
					{Group: tls.X25519},
				}},
				&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
					tls.PskModeDHE,
				}},
				&tls.SupportedVersionsExtension{Versions: []uint16{
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
			},
		}
	},
	Settings: []http2.Setting{
		{ID: http2.SettingInitialWindowSize, Val: 16777216},
	},
	ConnectionFlow: 16711681,
	PseudoHeaderOrder: []string{
		":method",
		":path",
		":authority",
		":scheme",
	},
}
