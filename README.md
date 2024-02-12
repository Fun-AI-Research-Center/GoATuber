<div align="center">
  <h1>Go-AI-vtuber</h1>
    <div>
    GoAtuber是一个基于golang、vue3、react和typescript的AI Vtber项目，通过简单的配置文件修改，即可打造出一个AI vtuber。
      <div><s>重构中。是的，这个项目还在重构。本次重构周期将会非常长。</s></div>
      <div>暂时暂停。主要维护者都在忙别的事了。恢复维护时间不定。但是肯定有更好用的框架，建议移步。这个项目更多是自娱自乐的产物，属于边写边学。</div>
    </div>
</div>

## 特色说明

- 写在前面：这篇`README.md`某些地方用词多少有点放荡不羁了，我觉得这也是一种特色。
- 最重要的：目前本项目的核心功能打包之后不需要手动安装任何依赖（除了你得拥有一个健全的操作系统和一个浏览器），开箱即用。就算你的电脑是10年前的老物（开发者使用ThinkPad 420完成了本项目的后端逻辑构建），只要能联网(你都到这里来了，联网是什么意思我不必解释了吧)，就能跑。
- 图形化的配置界面——还在完善之中。
- 同时支持VRM模型与Live2D模型。
- 支持GPT、Azure GPT、文心一言、讯飞星火等LLM平台，并且还将继续拓展。
- 支持使用Azure TTS、讯飞TTS进行文本转语音服务，并且还将继续拓展。
- 对句子进行划分，支持对每一句进行情感分析，让句子情感更细腻。
- 内置文本审查服务接口，可自由选择开启与否。
- 拥有直接对话接口，让皮套不止于vtuber——还可以是wife。
- 拥有直接朗读文本的接口，让你的AI vtuber不再仅仅是主播，也能是虚拟读稿皮套人。
- 拥有对话功能，支持使用Azure STT服务进行语音转文本。还将继续拓展。
- 未来将引入唱歌功能。
- 拥有在多平台进行直播的潜质。
- 众所周知（雾），项目开发者在B站吐槽了很多关于记忆的事情——在使用GPT的时候，本项目可以开启长期记忆模式。然而还是不完善，诚邀开发者一同解决。
- 不止于主播：有将项目引入QQ的计划。

- 也许是个特色？项目主创目前都刚大二。

## 使用指南

### 如何启动程序？

- 1.下载release中的压缩包，解压缩。
- 2.运行`GoATuber.exe`。
- 3.在浏览器弹出的窗口中选择并输入相应配置参数。如平台API Key等。
- 4.点击启动按钮，浏览器将弹出新的展示界面。

所以说，使用本项目的核心功能并不需要进行任何多余的配置操作。除了得获取各个平台的API Key等参数。

### 配置详解

教你如何获取各个配置参数，并填写相应配置。

#### Azure平台

不是打广告。申请github教育包，azure每年免费一百刀。足够玩了。

首先，在做事之前：

进入[Azure Portal](https://portal.azure.com)，创建一个[资源组](https://portal.azure.com/#view/HubsExtension/BrowseResourceGroups)。点击资源组页面的创建按钮，在创建页面的资源组输入框内给你的资源组取个名字。一路下一步，然后创建，完事。

- Azure OpenAI(GPT 服务等)
    - 1.在[Azure Portal](https://portal.azure.com)中，搜索并进入[Azure OpenAI](https://portal.azure.com/#blade/Microsoft_Azure_ProjectOxford/CognitiveServicesHub/OpenAI)。
    - 2.单击创建按钮。
    - 3.在创建界面里：
        - 1.选择之前创建的资源组。
        - 2.给服务起个名字。
        - 3.选择定价层。
    - 4.一路下一步(除非你有需求)，然后确认创建。
    - 5.在刚刚创建的服务中，选择下方的`Develop`块，可以看到两条秘钥(API Key)、你的region和终结点。
    - 6.将API Key(两个里选一个)和终结点复制填入`应用包配置-Azure`的相应输入框中。
    - 7.在`Get Started`块中，点击`Explore`按钮，进入Azure OPenAI工作室。
    - 8.在控制台左侧的导航栏中，选择`Deployments`选项。点击`Create new deployment`。在弹出的窗口中，第一个选择框选择使用的模型，第二个输入部署名称，作为deployment id(模型部署的名称)。
    - 9.创建资源。将模型部署的名称和所选择的模型填入配置文件界面。
    至此，Azure OpenAI的服务配置完成。
    - 10.关于温度等选项：日后再叙。

- Azure TTS服务与STT服务
    - 1.在[Azure Portal](https://portal.azure.com)中，搜索并进入[Speech service](https://portal.azure.com/#view/Microsoft_Azure_ProjectOxford/CognitiveServicesHub/~/SpeechServices)。
    - 2.服务创建流程参考Azure OpenAI。将配置参数填进配置文件栏里。
    - 3.TTS服务 打开[说话人列表](https://learn.microsoft.com/zh-cn/azure/ai-services/speech-service/language-support?tabs=tts#prebuilt-neural-voices)，挑选一个说话人，将文本转语音声音一栏的代号填入说话人名称中。
    - 4.TTS服务 语速、音量设置为基于基准量的正负百分比。请填写如:`+15.00%`的格式。
