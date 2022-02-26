package test

import (

	// "bitbucket.org/free5gc-team/openapi/models"

	"encoding/binary"
	"fmt"
	"net"
	"test/ngapTestpacket"

	"bitbucket.org/free5gc-team/nas"
	"bitbucket.org/free5gc-team/nas/nasMessage"
	"bitbucket.org/free5gc-team/ngap"
)

// This function is used for nas packet
func DecodePDUSessionEstablishmentAccept(ue *RanUeContext, length int, buffer []byte) (*nas.Message, error) {

	if length == 0 {
		return nil, fmt.Errorf("Empty buffer")
	}

	nasEnv, n := DecapNasPduFromEnvelope(buffer[:length])
	nasMsg, err := NASDecode(ue, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, nasEnv[:n])
	if err != nil {
		return nil, fmt.Errorf("NAS Decode Fail: %+v", err)
	}

	// Retrieve GSM from GmmMessage.DLNASTransport.PayloadContainer and decode
	payloadContainer := nasMsg.GmmMessage.DLNASTransport.PayloadContainer
	byteArray := payloadContainer.Buffer[:payloadContainer.Len]
	if err := nasMsg.GsmMessageDecode(&byteArray); err != nil {
		return nil, fmt.Errorf("NAS Decode Fail: %+v", err)
	}

	return nasMsg, nil
}

// This function is used for nas packet
func GetPDUAddress(accept *nasMessage.PDUSessionEstablishmentAccept) (net.IP, error) {
	if addr := accept.PDUAddress; addr != nil {
		PDUSessionTypeValue := addr.GetPDUSessionTypeValue()
		if PDUSessionTypeValue == nasMessage.PDUSessionTypeIPv4 {
			ip := net.IP(addr.Octet[1:5])
			return ip, nil
		}
	}

	return nil, fmt.Errorf("PDUAddress is nil")
}

func GetNGSetupRequest(gnbId []byte, bitlength uint64, name string) ([]byte, error) {
	message := ngapTestpacket.BuildNGSetupRequest()
	// GlobalRANNodeID
	ie := message.InitiatingMessage.Value.NGSetupRequest.ProtocolIEs.List[0]
	gnbID := ie.Value.GlobalRANNodeID.GlobalGNBID.GNBID.GNBID
	gnbID.Bytes = gnbId
	gnbID.BitLength = bitlength
	// RANNodeName
	ie = message.InitiatingMessage.Value.NGSetupRequest.ProtocolIEs.List[1]
	ie.Value.RANNodeName.Value = name

	return ngap.Encoder(message)
}

func GetInitialUEMessage(ranUeNgapID int64, nasPdu []byte, fiveGSTmsi string) ([]byte, error) {
	message := ngapTestpacket.BuildInitialUEMessage(ranUeNgapID, nasPdu, fiveGSTmsi)
	return ngap.Encoder(message)
}

func GetUplinkNASTransport(amfUeNgapID, ranUeNgapID int64, nasPdu []byte) ([]byte, error) {
	message := ngapTestpacket.BuildUplinkNasTransport(amfUeNgapID, ranUeNgapID, nasPdu)
	return ngap.Encoder(message)
}

func GetInitialContextSetupResponse(amfUeNgapID int64, ranUeNgapID int64) ([]byte, error) {
	message := ngapTestpacket.BuildInitialContextSetupResponseForRegistraionTest(amfUeNgapID, ranUeNgapID)

	return ngap.Encoder(message)
}

func GetInitialContextSetupResponseForServiceRequest(
	amfUeNgapID int64, ranUeNgapID int64, ipv4 string) ([]byte, error) {
	message := ngapTestpacket.BuildInitialContextSetupResponse(amfUeNgapID, ranUeNgapID, ipv4, nil)
	return ngap.Encoder(message)
}

func GetPDUSessionResourceSetupResponse(pduSessionId int64, amfUeNgapID int64, ranUeNgapID int64, ipv4 string) ([]byte, error) {
	message := ngapTestpacket.BuildPDUSessionResourceSetupResponseForRegistrationTest(pduSessionId, amfUeNgapID, ranUeNgapID, ipv4)
	return ngap.Encoder(message)
}

func EncodeNasPduWithSecurity(ue *RanUeContext, pdu []byte, securityHeaderType uint8,
	securityContextAvailable, newSecurityContext bool) ([]byte, error) {
	m := nas.NewMessage()
	err := m.PlainNasDecode(&pdu)
	if err != nil {
		return nil, err
	}
	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    securityHeaderType,
	}
	return NASEncode(ue, m, securityContextAvailable, newSecurityContext)

}

func EncodeNasPduInEnvelopeWithSecurity(ue *RanUeContext, pdu []byte, securityHeaderType uint8,
	securityContextAvailable, newSecurityContext bool) ([]byte, error) {
	m := nas.NewMessage()
	err := m.PlainNasDecode(&pdu)
	if err != nil {
		return nil, err
	}
	m.SecurityHeader = nas.SecurityHeader{
		ProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		SecurityHeaderType:    securityHeaderType,
	}
	return NASEnvelopeEncode(ue, m, securityContextAvailable, newSecurityContext)

}

func DecapNasPduFromEnvelope(envelop []byte) ([]byte, int) {
	// According to TS 24.502 8.2.4 and TS 24.502 9.4,
	// a NAS message envelope = Length | NAS Message

	// Get NAS Message Length
	nasLen := binary.BigEndian.Uint16(envelop[:2])
	nasMsg := make([]byte, nasLen)
	copy(nasMsg, envelop[2:2+nasLen])

	return nasMsg, int(nasLen)
}

func GetUEContextReleaseComplete(amfUeNgapID int64, ranUeNgapID int64, pduSessionIDList []int64) ([]byte, error) {
	message := ngapTestpacket.BuildUEContextReleaseComplete(amfUeNgapID, ranUeNgapID, pduSessionIDList)
	return ngap.Encoder(message)
}

func GetUEContextReleaseRequest(amfUeNgapID int64, ranUeNgapID int64, pduSessionIDList []int64) ([]byte, error) {
	message := ngapTestpacket.BuildUEContextReleaseRequest(amfUeNgapID, ranUeNgapID, pduSessionIDList)
	return ngap.Encoder(message)
}

func GetPDUSessionResourceReleaseResponse(amfUeNgapID int64, ranUeNgapID int64) ([]byte, error) {
	message := ngapTestpacket.BuildPDUSessionResourceReleaseResponseForReleaseTest(amfUeNgapID, ranUeNgapID)
	return ngap.Encoder(message)
}
func GetPathSwitchRequest(amfUeNgapID int64, ranUeNgapID int64) ([]byte, error) {
	message := ngapTestpacket.BuildPathSwitchRequest(amfUeNgapID, ranUeNgapID)
	message.InitiatingMessage.Value.PathSwitchRequest.ProtocolIEs.List =
		message.InitiatingMessage.Value.PathSwitchRequest.ProtocolIEs.List[0:5]
	return ngap.Encoder(message)
}

func GetHandoverRequired(
	amfUeNgapID int64, ranUeNgapID int64, targetGNBID []byte, targetCellID []byte) ([]byte, error) {
	message := ngapTestpacket.BuildHandoverRequired(amfUeNgapID, ranUeNgapID, targetGNBID, targetCellID)
	return ngap.Encoder(message)
}

func GetHandoverRequestAcknowledge(amfUeNgapID int64, ranUeNgapID int64) ([]byte, error) {
	message := ngapTestpacket.BuildHandoverRequestAcknowledge(amfUeNgapID, ranUeNgapID)
	return ngap.Encoder(message)
}

func GetHandoverNotify(amfUeNgapID int64, ranUeNgapID int64) ([]byte, error) {
	message := ngapTestpacket.BuildHandoverNotify(amfUeNgapID, ranUeNgapID)
	return ngap.Encoder(message)
}

func GetPDUSessionResourceSetupResponseForPaging(amfUeNgapID int64, ranUeNgapID int64, ipv4 string) ([]byte, error) {
	message := ngapTestpacket.BuildPDUSessionResourceSetupResponseForPaging(amfUeNgapID, ranUeNgapID, ipv4)
	return ngap.Encoder(message)
}
