const fs = require("fs");
const util = require("util");

const path = require("path");

const statAsync = util.promisify(fs.stat);
const readdirAsync = util.promisify(fs.readdir);

async function getFiles(dir) {
    const subdirs = await readdirAsync(dir);
    const files = await Promise.all(
        subdirs.map(async (subdir) => {
            const res = path.resolve(dir, subdir);
            return (await statAsync(res)).isDirectory() ? getFiles(res) : res;
        })
    );

    return files.reduce((a, f) => a.concat(f), []);
}
