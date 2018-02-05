function getNewKey() {
  var ret = new Uint8Array(256 / 8);
  window.crypto.getRandomValues(ret);
  return ret;
}

function encryptMsg(cleartext, key) {
  var ret = "";
  var text_bytes = aesjs.utils.utf8.toBytes(cleartext);
  var aes_ctr = new aesjs.ModeOfOperation.ctr(key, new aesjs.Counter(5));
  var encrypted_bytes = aes_ctr.encrypt(text_bytes);
  ret = aesjs.utils.hex.fromBytes(encrypted_bytes);
  return ret;
}

function encryptPaste() {
  var ret = true;
  var key = getNewKey();
  var cleartext = document.paste_form.paste_value.value;
  var cyphertext = encryptMsg(cleartext, key);
  var hex_key = aesjs.utils.hex.fromBytes(key);
  document.paste_form.paste_key.value = hex_key;
  document.paste_form.paste_value.value = cyphertext;
  return true;
}
