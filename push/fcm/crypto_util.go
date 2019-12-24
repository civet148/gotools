// Copyright (c) 2019-present,  NebulaChat Studio (https://nebula.chat).
//  All rights reserved.
//
// Author: Benqi (wubenqi@gmail.com)
//

package fcm

import (
	"nebula.chat/enterprise/mtproto"
	"nebula.chat/enterprise/pkg/hack"
)

/**
{
	Object value = data.get("p");
	if (!(value instanceof String)) {
		if (BuildVars.LOGS_ENABLED) {
			FileLog.d("GCM DECRYPT ERROR 1");
		}
		onDecryptError();
		return;
	}
	byte[] bytes = Base64.decode((String) value, Base64.URL_SAFE);
	NativeByteBuffer buffer = new NativeByteBuffer(bytes.length);
	buffer.writeBytes(bytes);
	buffer.position(0);

	if (SharedConfig.pushAuthKeyId == null) {
		SharedConfig.pushAuthKeyId = new byte[8];
		byte[] authKeyHash = Utilities.computeSHA1(SharedConfig.pushAuthKey);
		System.arraycopy(authKeyHash, authKeyHash.length - 8, SharedConfig.pushAuthKeyId, 0, 8);
	}
	byte[] inAuthKeyId = new byte[8];
	buffer.readBytes(inAuthKeyId, true);
	if (!Arrays.equals(SharedConfig.pushAuthKeyId, inAuthKeyId)) {
		onDecryptError();
		if (BuildVars.LOGS_ENABLED) {
			FileLog.d(String.format(Locale.US, "GCM DECRYPT ERROR 2 k1=%s k2=%s, key=%s", Utilities.bytesToHex(SharedConfig.pushAuthKeyId), Utilities.bytesToHex(inAuthKeyId), Utilities.bytesToHex(SharedConfig.pushAuthKey)));
		}
		return;
	}

	byte[] messageKey = new byte[16];
	buffer.readBytes(messageKey, true);

	MessageKeyData messageKeyData = MessageKeyData.generateMessageKeyData(SharedConfig.pushAuthKey, messageKey, true, 2);
	Utilities.aesIgeEncryption(buffer.buffer, messageKeyData.aesKey, messageKeyData.aesIv, false, false, 24, bytes.length - 24);

	byte[] messageKeyFull = Utilities.computeSHA256(SharedConfig.pushAuthKey, 88 + 8, 32, buffer.buffer, 24, buffer.buffer.limit());
	if (!Utilities.arraysEquals(messageKey, 0, messageKeyFull, 8)) {
		onDecryptError();
		if (BuildVars.LOGS_ENABLED) {
			FileLog.d(String.format("GCM DECRYPT ERROR 3, key = %s", Utilities.bytesToHex(SharedConfig.pushAuthKey)));
		}
		return;
	}

	int len = buffer.readInt32(true);
	byte[] strBytes = new byte[len];
	buffer.readBytes(strBytes, true);
	jsonString = new String(strBytes, "UTF-8");
	JSONObject json = new JSONObject(jsonString);
}
*/

func cryptoPushData(jsonString string, pushAuthKey []byte) ([]byte, error) {
	b := hack.Bytes(jsonString)
	x := mtproto.NewEncodeBuf(len(b) + 4)
	x.Int(int32(len(b)))
	x.Bytes(b)

	return []byte{}, nil
}
