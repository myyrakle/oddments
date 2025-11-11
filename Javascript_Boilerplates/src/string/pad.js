function leftPad(len, character) {
    var self = String(this);
    len = len - self.length;
    if (len <= 0) {
        return self;
    }
    character = String(character);

    var pad = "";
    while (true) {
        if (len & 1) {
            pad += character;
        }
        len >>= 1;
        if (len) {
            character += character;
        } else {
            break;
        }
    }

    return pad + self;
};

function rightPad(len, character) {
    var self = String(this);
    len = len - self.length;
    if (len <= 0) {
        return self;
    }
    character = String(character);

    var pad = "";
    while (true) {
        if (len & 1) {
            pad += character;
        }
        len >>= 1;
        if (len) {
            character += character;
        } else {
            break;
        }
    }

    return self + pad;
};
