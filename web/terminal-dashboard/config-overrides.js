const { override, addBabelPlugin, addWebpackPlugin } = require('customize-cra');
const JavaScriptObfuscator = require('webpack-obfuscator');

// 自定义 Webpack 插件封装
class JavaScriptObfuscatorPlugin {
    constructor(options) {
        this.options = options;
    }

    apply(compiler) {
        new JavaScriptObfuscator(this.options).apply(compiler);
    }
}

module.exports = override(
    process.env.NODE_ENV === 'production' && addBabelPlugin('transform-remove-console'),
    process.env.NODE_ENV === 'production' && addWebpackPlugin(new JavaScriptObfuscatorPlugin({
        rotateUnicodeArray: true,
        compact: true,
        controlFlowFlattening: true,
        deadCodeInjection: true,
        debugProtection: true,
        disableConsoleOutput: true,
        identifierNamesGenerator: 'hexadecimal',
        selfDefending: true,
        stringArray: true,
        stringArrayEncoding: ['base64'],
        stringArrayThreshold: 0.75
    }))
);
