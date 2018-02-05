function getKeyFromUri() {
  var ret = "";
  var hex_key = window.location.hash.slice(1);
  ret = aesjs.utils.hex.toBytes(hex_key);
  console.log(ret);
  return ret;
}

function decryptMsg(cyphertext, key) {
  var ret = "";
  var encrypted_bytes = aesjs.utils.hex.toBytes(cyphertext);
  var aes_ctr = new aesjs.ModeOfOperation.ctr(key, new aesjs.Counter(5));
  var cleartext_bytes = aes_ctr.decrypt(encrypted_bytes);
  ret = aesjs.utils.utf8.fromBytes(cleartext_bytes);
  return ret;
}

$(document).ready(function() {
  var key = getKeyFromUri();
  var msg_container = document.getElementById("cyphertext");
  var cyphertext = msg_container.innerHTML;
  var cleartext = decryptMsg(cyphertext, key);
  cleartext = cleartext.replace(/\n/g, "<br/>");
  msg_container.innerHTML = cleartext;
});
