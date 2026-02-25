package profiles

import (
	"time"

	"github.com/nukilabs/http"
	"github.com/nukilabs/http/http2"
	"github.com/nukilabs/quic-go/http3"
	tls "github.com/nukilabs/utls"
)

func chromeweight(p http.Priority) uint8 {
	switch p {
	case http.PriorityHighest:
		return 255
	case http.PriorityHigh:
		return 219
	case http.PriorityMedium:
		return 183
	case http.PriorityLow:
		return 146
	case http.PriorityLowest:
		return 109
	case http.PriorityIdle:
		return 73
	case http.PriorityThrottled:
		return 36
	default:
		return 219
	}
}

func Chrome(major int) ClientProfile {
	switch major {
	case 120:
		return Chrome120
	case 121:
		return Chrome121
	case 122:
		return Chrome122
	case 123:
		return Chrome123
	case 124:
		return Chrome124
	case 125:
		return Chrome125
	case 126:
		return Chrome126
	case 127:
		return Chrome127
	case 128:
		return Chrome128
	case 129:
		return Chrome129
	case 130:
		return Chrome130
	case 131:
		return Chrome131
	case 132:
		return Chrome132
	case 133:
		return Chrome133
	case 134:
		return Chrome134
	case 135:
		return Chrome135
	case 136:
		return Chrome136
	case 137:
		return Chrome137
	case 138:
		return Chrome138
	case 139:
		return Chrome139
	case 140:
		return Chrome140
	case 141:
		return Chrome141
	case 142:
		return Chrome142
	case 143:
		return Chrome143
	case 144:
		return Chrome144
	case 145:
		return Chrome145
	default:
		return Chrome133
	}
}

var Chrome120 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.GREASE_PLACEHOLDER,
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []byte{
				0x00, // compressionNone
			},
			Extensions: tls.ShuffleChromeTLSExtensions([]tls.TLSExtension{
				&tls.UtlsGREASEExtension{},
				&tls.SNIExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.GREASE_PLACEHOLDER,
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00, // pointFormatUncompressed
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
				}},
				&tls.SCTExtension{},
				&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
					{Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
					{Group: tls.X25519},
				}},
				&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
					tls.PskModeDHE,
				}},
				&tls.SupportedVersionsExtension{Versions: []uint16{
					tls.GREASE_PLACEHOLDER,
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.UtlsCompressCertExtension{Algorithms: []tls.CertCompressionAlgo{
					tls.CertCompressionBrotli,
				}},
				&tls.ApplicationSettingsExtension{
					SupportedProtocols: []string{"h2"},
				},
				tls.BoringGREASEECH(),
				&tls.UtlsGREASEExtension{},
				&tls.UtlsPreSharedKeyExtension{},
			}),
		}
	},
	H2: &H2ClientProfile{
		Settings: []http2.Setting{
			{ID: http2.SettingHeaderTableSize, Val: 65536},
			{ID: http2.SettingEnablePush, Val: 0},
			{ID: http2.SettingInitialWindowSize, Val: 6291456},
			{ID: http2.SettingMaxHeaderListSize, Val: 262144},
		},
		ConnectionFlow: 15663105,
		HeaderPriority: func(r *http.Request) http2.PriorityParam {
			return http2.PriorityParam{
				StreamDep: 0,
				Exclusive: true,
				Weight:    chromeweight(r.Priority),
			}
		},
	},
	PseudoHeaderOrder: []string{
		":method",
		":authority",
		":scheme",
		":path",
	},
}

var Chrome121 = Chrome120

var Chrome122 = Chrome120

var Chrome123 = Chrome120

var Chrome124 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.GREASE_PLACEHOLDER,
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []byte{
				0x00, // compressionNone
			},
			Extensions: tls.ShuffleChromeTLSExtensions([]tls.TLSExtension{
				&tls.UtlsGREASEExtension{},
				&tls.SNIExtension{},
				&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
					{Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
					{Group: tls.X25519Kyber768Draft00},
					{Group: tls.X25519},
				}},
				&tls.UtlsCompressCertExtension{Algorithms: []tls.CertCompressionAlgo{
					tls.CertCompressionBrotli,
				}},
				&tls.SCTExtension{},
				&tls.SupportedVersionsExtension{Versions: []uint16{
					tls.GREASE_PLACEHOLDER,
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.StatusRequestExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.SessionTicketExtension{},
				&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
					tls.PskModeDHE,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.CurveID(tls.GREASE_PLACEHOLDER),
					tls.X25519Kyber768Draft00,
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00, // pointFormatUncompressed
				}},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PSSWithSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PSSWithSHA384,
					tls.PKCS1WithSHA384,
					tls.PSSWithSHA512,
					tls.PKCS1WithSHA512,
				}},
				&tls.ApplicationSettingsExtension{
					SupportedProtocols: []string{"h2"},
				},
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateOnceAsClient},
				&tls.ExtendedMasterSecretExtension{},
				tls.BoringGREASEECH(),
				&tls.UtlsGREASEExtension{},
				&tls.UtlsPreSharedKeyExtension{},
			}),
		}
	},
	H2: &H2ClientProfile{
		Settings: []http2.Setting{
			{ID: http2.SettingHeaderTableSize, Val: 65536},
			{ID: http2.SettingEnablePush, Val: 0},
			{ID: http2.SettingInitialWindowSize, Val: 6291456},
			{ID: http2.SettingMaxHeaderListSize, Val: 262144},
		},
		ConnectionFlow: 15663105,
		HeaderPriority: func(r *http.Request) http2.PriorityParam {
			return http2.PriorityParam{
				StreamDep: 0,
				Exclusive: true,
				Weight:    chromeweight(r.Priority),
			}
		},
	},
	PseudoHeaderOrder: []string{
		":method",
		":authority",
		":scheme",
		":path",
	},
}

var Chrome125 = Chrome124

var Chrome126 = Chrome124

var Chrome127 = Chrome124

var Chrome128 = Chrome124

var Chrome129 = Chrome124

var Chrome130 = Chrome124

var Chrome131 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.GREASE_PLACEHOLDER,
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []byte{
				0x00, // compressionNone
			},
			Extensions: tls.ShuffleChromeTLSExtensions([]tls.TLSExtension{
				&tls.UtlsGREASEExtension{},
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateNever},
				&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
					tls.PskModeDHE,
				}},
				&tls.SessionTicketExtension{},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{0x00}},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PSSWithSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PSSWithSHA384,
					tls.PKCS1WithSHA384,
					tls.PSSWithSHA512,
					tls.PKCS1WithSHA512,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.CurveID(tls.GREASE_PLACEHOLDER),
					tls.X25519MLKEM768,
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
				&tls.SNIExtension{},
				tls.BoringGREASEECH(),
				&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
					{Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
					{Group: tls.X25519MLKEM768},
					{Group: tls.X25519},
				}},
				&tls.UtlsCompressCertExtension{Algorithms: []tls.CertCompressionAlgo{
					tls.CertCompressionBrotli,
				}},
				&tls.ApplicationSettingsExtension{
					SupportedProtocols: []string{"h2"},
				},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.StatusRequestExtension{},
				&tls.SupportedVersionsExtension{Versions: []uint16{
					tls.GREASE_PLACEHOLDER,
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.SCTExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.UtlsGREASEExtension{},
				&tls.UtlsPreSharedKeyExtension{OmitEmptyPsk: true},
			}),
		}
	},
	H2: &H2ClientProfile{
		Settings: []http2.Setting{
			{ID: http2.SettingHeaderTableSize, Val: 65536},
			{ID: http2.SettingEnablePush, Val: 0},
			{ID: http2.SettingInitialWindowSize, Val: 6291456},
			{ID: http2.SettingMaxHeaderListSize, Val: 262144},
		},
		ConnectionFlow: 15663105,
		HeaderPriority: func(r *http.Request) http2.PriorityParam {
			return http2.PriorityParam{
				StreamDep: 0,
				Exclusive: true,
				Weight:    chromeweight(r.Priority),
			}
		},
	},
	PseudoHeaderOrder: []string{
		":method",
		":authority",
		":scheme",
		":path",
	},
}

var Chrome132 = Chrome131

var Chrome133 = ClientProfile{
	ClientHelloSpec: func() *tls.ClientHelloSpec {
		return &tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.GREASE_PLACEHOLDER,
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []byte{
				0x00, // compressionNone
			},
			Extensions: tls.ShuffleChromeTLSExtensions([]tls.TLSExtension{
				&tls.UtlsGREASEExtension{},
				&tls.RenegotiationInfoExtension{Renegotiation: tls.RenegotiateNever},
				&tls.PSKKeyExchangeModesExtension{Modes: []uint8{
					tls.PskModeDHE,
				}},
				&tls.SessionTicketExtension{},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{0x00}},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []tls.SignatureScheme{
					tls.ECDSAWithP256AndSHA256,
					tls.PSSWithSHA256,
					tls.PKCS1WithSHA256,
					tls.ECDSAWithP384AndSHA384,
					tls.PSSWithSHA384,
					tls.PKCS1WithSHA384,
					tls.PSSWithSHA512,
					tls.PKCS1WithSHA512,
				}},
				&tls.SupportedCurvesExtension{Curves: []tls.CurveID{
					tls.CurveID(tls.GREASE_PLACEHOLDER),
					tls.X25519MLKEM768,
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
				&tls.SNIExtension{},
				tls.BoringGREASEECH(),
				&tls.KeyShareExtension{KeyShares: []tls.KeyShare{
					{Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
					{Group: tls.X25519MLKEM768},
					{Group: tls.X25519},
				}},
				&tls.UtlsCompressCertExtension{Algorithms: []tls.CertCompressionAlgo{
					tls.CertCompressionBrotli,
				}},
				&tls.ApplicationSettingsExtensionNew{
					SupportedProtocols: []string{"h2"},
				},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.StatusRequestExtension{},
				&tls.SupportedVersionsExtension{Versions: []uint16{
					tls.GREASE_PLACEHOLDER,
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.SCTExtension{},
				&tls.ExtendedMasterSecretExtension{},
				&tls.UtlsGREASEExtension{},
				&tls.UtlsPreSharedKeyExtension{OmitEmptyPsk: true},
			}),
		}
	},
	H2: &H2ClientProfile{
		Settings: []http2.Setting{
			{ID: http2.SettingHeaderTableSize, Val: 65536},
			{ID: http2.SettingEnablePush, Val: 0},
			{ID: http2.SettingInitialWindowSize, Val: 6291456},
			{ID: http2.SettingMaxHeaderListSize, Val: 262144},
		},
		ConnectionFlow: 15663105,
		HeaderPriority: func(r *http.Request) http2.PriorityParam {
			return http2.PriorityParam{
				StreamDep: 0,
				Exclusive: true,
				Weight:    chromeweight(r.Priority),
			}
		},
		InflowTimeout:   5 * time.Second,
		ReadIdleTimeout: 10 * time.Second,
		PrefacePing:     true,
	},
	H3: &H3ClientProfile{
		Settings: []http3.Setting{
			{ID: http3.SettingQpackMaxTableCapacity, Val: 65536},
			{ID: http3.SettingMaxFieldSectionSize, Val: 262144},
			{ID: http3.SettingQpackBlockedStreams, Val: 100},
			{ID: http3.SettingH3Datagram, Val: 1},
			{ID: http3.SettingGrease, Val: 0},
		},
	},
	PseudoHeaderOrder: []string{
		":method",
		":authority",
		":scheme",
		":path",
	},
}

var Chrome134 = Chrome133

var Chrome135 = Chrome133

var Chrome136 = Chrome133

var Chrome137 = Chrome133

var Chrome138 = Chrome133

var Chrome139 = Chrome133

var Chrome140 = Chrome133

var Chrome141 = Chrome133

var Chrome142 = Chrome133

var Chrome143 = Chrome133

var Chrome144 = Chrome133

var Chrome145 = Chrome133
