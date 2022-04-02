package ws_stream

import (
	"encoding/binary"
	"errors"
	"games/comm/utils"
	"games/core/conn/def"
	"games/core/conn/transmit"
	"log"

	"github.com/gorilla/websocket"
)

/// <summary>
/// Channel websocket协议读写解析
/// <summary>
type Channel struct {
}

func NewChannel() transmit.IChannel {
	return &Channel{}
}

func (s *Channel) OnRecv(conn interface{}) (interface{}, error, def.Reason) {
	c, ok := conn.(*websocket.Conn)
	if !ok || c == nil {
		panic(errors.New("ws_stream.OnRecv conn == nil"))
	}
	//len+CRC，4字节
	//conn.SetReadLimit(4)
	msgType, buf, err := c.ReadMessage()
	if err != nil {
		return nil, err, def.KClosed
	}
	//TextMessage/BinaryMessage
	if websocket.BinaryMessage != msgType {
		log.Fatalln("OnRecvMessage: msgType error")
		return nil, errors.New("ws_stream.OnRecv parse error"), def.KExcept
	}
	//len，2字节
	length := binary.LittleEndian.Uint16(buf[:2])
	if length != uint16(len(buf)) {
		log.Fatalln("OnRecvMessage: checklen error")
		return nil, errors.New("ws_stream.OnRecv parse error"), def.KExcept
	}
	return buf, nil, def.KNoError
	//CRC，2字节
	//chsum := binary.LittleEndian.Uint16(buf[2:])
	// 读取剩余大小
	//conn.SetReadLimit(int64(len - 4))
	//_, remain, err := conn.ReadMessage()
	//if err != nil {
	//	log.Fatalln("OnRecvMessage: ", err)
	//	return nil, err
	//}
	//CRC校验
	// crc := GetChecksum(buf[4:])
	// if crc != chsum {
	// 	log.Fatalln("OnRecvMessage: GetChecksum error")
	// 	return nil, errors.New("GetChecksum error")
	// }
	// //版本0x0001
	// ver := binary.LittleEndian.Uint16(buf[4:])
	// //标记0x5F5F
	// sign := binary.LittleEndian.Uint16(buf[6:])
	// //主命令ID
	//mainID := uint8(buf[8])
	//子命令ID
	//subID := uint8(buf[9])
	// //加密类型
	// encryptTy := uint8(buf[10])
	// //预留字段
	// reserve := uint8(buf[11])
	// //请求ID
	// reqID := binary.LittleEndian.Uint32(buf[12:16])
	// //实际大小
	// realSize := binary.LittleEndian.Uint16(buf[16:18])
	// log.Printf("ver:%#x\nsign:%#x\nmainID:%d\nsubID:%d\nencTy:%#x\nreserv:%d\nreqID:%d\nrealSize:%d\n",
	// 	ver, sign, mainID, subID, encryptTy, reserve, reqID, realSize)
	//实际protobuf数据
	//data := buf[18:]
	//msg, _, err = codec.DecodeMessage(int(subID), data)
	//if err != nil {
	//	log.Println("MyWSChannel::OnSendMessage ", msg)
	//}
	//便于框架处理
	// msg := &RootMsg{}
	// msg.Cmd = uint32(ENWORD(int(mainID), int(subID)))
	// if mainID == uint8(Game_Common.MAINID_MAIN_MESSAGE_CLIENT_TO_GAME_LOGIC) {
	// 	tMsg, _, err := codec.DecodeMessage(SubCmdID, buf[18:])
	// 	if err != nil {
	// 		log.Printf("OnRecvMessage 1: [mainID=%d subID=%d] ERR: %v\n", mainID, subID, err)
	// 		return nil, err
	// 	}
	// 	pMsg, ok := tMsg.(*GameServer.MSG_CSC_Passageway)
	// 	if !ok {
	// 		log.Printf("OnRecvMessage 2: [mainID=%d subID=%d] ERR: %v\n", mainID, subID, err)
	// 		return nil, nil
	// 	}
	// 	log.Printf("OnRecvMessage 0: [mainID=%d subID=%d]\n%v\n%v\n", mainID, subID, reflect.TypeOf(pMsg).Elem(), util.JSON2Str(pMsg))
	// 	//log.Printf("PassData.len = %d\n", len(pMsg.PassData[:]))
	// 	//log.Printf("\n\n%v\n\n", pMsg.PassData[:])
	// 	msg.Data = pMsg.PassData[:]
	// } else {
	// 	msg.Data = buf[18:]
	// }
}

func (s *Channel) OnSend(conn interface{}, msg interface{}) (error, def.Reason) {
	c, ok := conn.(*websocket.Conn)
	if !ok || c == nil {
		panic(errors.New("ws_stream.OnSend conn == nil"))
	}
	buf, _ := utils.ToBytes(msg)
	return c.WriteMessage(websocket.BinaryMessage, buf), def.KNoError
	//log.Println("MyWSChannel::OnSendMessage ", msg)
	// h, ok := msg.(*Msg)
	// if !ok || h == nil {
	// 	return nil
	// }
	// data, _, err := codec.EncodeMessage(h.msg, nil)
	// if err != nil {
	// 	log.Fatalln("EncodeMessage : ", err)
	// 	return err
	// }
	// buf := make([]byte, 18+len(data))
	// //len，2字节
	// length := 18 + len(data)
	// binary.LittleEndian.PutUint16(buf[0:], uint16(length))
	// //CRC，2字节
	// //binary.LittleEndian.PutUint16(buf[2:], h.crc)
	// //版本0x0001
	// binary.LittleEndian.PutUint16(buf[4:], uint16(h.ver))
	// //标记0x5F5F
	// binary.LittleEndian.PutUint16(buf[6:], uint16(h.sign))
	// //主命令ID
	// buf[8] = byte(h.mainID)
	// //子命令ID
	// buf[9] = byte(h.subID)
	// //加密类型
	// buf[10] = byte(h.encType)
	// //预留字段
	// buf[11] = byte(0x01) //
	// //请求ID
	// binary.LittleEndian.PutUint32(buf[12:], uint32(0)) //
	// //实际大小
	// binary.LittleEndian.PutUint16(buf[16:], uint16(len(data)))
	// //实际数据
	// copy(buf[18:], data)
	// //CRC，2字节
	// crc := GetChecksum(buf[4:])
	// binary.LittleEndian.PutUint16(buf[2:], crc)
	// return conn.WriteMessage(websocket.BinaryMessage, buf)
}
