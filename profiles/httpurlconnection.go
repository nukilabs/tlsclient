package profiles

import (
	tls "github.com/nukilabs/utls"
)

func HttpUrlConnectionAndroid(sdk int) ClientProfile {
	switch sdk {
	case 21:
		return HttpUrlConnectionAndroid21
	case 22:
		return HttpUrlConnectionAndroid22
	case 23:
		return HttpUrlConnectionAndroid23
	case 24:
		return HttpUrlConnectionAndroid24
	case 25:
		return HttpUrlConnectionAndroid25
	case 26:
		return HttpUrlConnectionAndroid26
	case 27:
		return HttpUrlConnectionAndroid27
	case 28:
		return HttpUrlConnectionAndroid28
	case 29:
		return HttpUrlConnectionAndroid29
	case 30:
		return HttpUrlConnectionAndroid30
	case 31:
		return HttpUrlConnectionAndroid31
	case 32:
		return HttpUrlConnectionAndroid32
	case 33:
		return HttpUrlConnectionAndroid33
	case 34:
		return HttpUrlConnectionAndroid34
	case 35:
		return HttpUrlConnectionAndroid35
	default:
		return HttpUrlConnectionAndroid29
	}
}

var HttpUrlConnectionAndroid21 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_128_CBC_SHA,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_256_CBC_SHA,
				tls.FAKE_TLS_DHE_DSS_WITH_AES_128_CBC_SHA,
				uint16(0x0038),
				tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
				tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_RC4_128_SHA,
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
			},
		}
	},
}

var HttpUrlConnectionAndroid22 = HttpUrlConnectionAndroid21

var HttpUrlConnectionAndroid23 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_128_CBC_SHA,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
				tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_RC4_128_SHA,
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
				&tls.ALPNExtension{AlpnProtocols: []string{"http/1.1"}},
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
}

var HttpUrlConnectionAndroid24 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_128_CBC_SHA,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_256_CBC_SHA,
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
				&tls.ALPNExtension{AlpnProtocols: []string{"http/1.1"}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.CurveP256,
				}},
			},
		}
	},
}

var HttpUrlConnectionAndroid25 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_128_CBC_SHA,
				tls.FAKE_TLS_DHE_RSA_WITH_AES_256_CBC_SHA,
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
				&tls.ALPNExtension{AlpnProtocols: []string{"http/1.1"}},
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
}

var HttpUrlConnectionAndroid26 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
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
				&tls.ALPNExtension{AlpnProtocols: []string{"http/1.1"}},
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
}

var HttpUrlConnectionAndroid27 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
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
				&tls.ALPNExtension{AlpnProtocols: []string{"http/1.1"}},
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
}

var HttpUrlConnectionAndroid28 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
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
				&tls.ALPNExtension{AlpnProtocols: []string{"http/1.1"}},
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
}

var HttpUrlConnectionAndroid29 = ClientProfile{
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
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
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
				&tls.ALPNExtension{AlpnProtocols: []string{"http/1.1"}},
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
					tls.VersionTLS11,
					tls.VersionTLS10,
				}},
				&tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
			},
		}
	},
}

var HttpUrlConnectionAndroid30 = HttpUrlConnectionAndroid29

var HttpUrlConnectionAndroid31 = HttpUrlConnectionAndroid29

var HttpUrlConnectionAndroid32 = HttpUrlConnectionAndroid29

var HttpUrlConnectionAndroid33 = HttpUrlConnectionAndroid29

var HttpUrlConnectionAndroid34 = HttpUrlConnectionAndroid29

var HttpUrlConnectionAndroid35 = HttpUrlConnectionAndroid29
