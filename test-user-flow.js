/**
 * B站商品评论分析工具 - 完整用户流程测试
 * 
 * 测试目标：
 * 1. 验证首页自由需求输入功能
 * 2. 验证快捷示例是否为需求描述
 * 3. 验证输入测试需求并提交
 * 4. 验证跳转到确认页
 * 5. 验证确认页显示所有必要信息
 * 6. 截图保存确认页面
 */

const { chromium } = require('playwright');

(async () => {
  console.log('🚀 开始测试 B站商品评论分析工具用户流程...\n');
  
  // 启动浏览器（使用系统已安装的浏览器）
  const browser = await chromium.launch({ 
    headless: true,   // 无头模式，不显示浏览器窗口
    channel: 'chrome' // 使用系统Chrome浏览器
  });
  
  const context = await browser.newContext({
    viewport: { width: 1280, height: 720 }
  });
  
  const page = await context.newPage();
  
  try {
    // ========== 步骤1: 访问首页 ==========
    console.log('✅ 步骤1: 访问首页 http://localhost:5173');
    await page.goto('http://localhost:5173', { waitUntil: 'networkidle' });
    await page.waitForTimeout(1000);
    
    // ========== 步骤2: 检查输入框是否为自由需求输入 ==========
    console.log('✅ 步骤2: 检查输入框是否为自由需求输入');
    const inputPlaceholder = await page.locator('input[type="text"], textarea').first().getAttribute('placeholder');
    console.log(`   输入框提示文字: "${inputPlaceholder}"`);
    
    if (inputPlaceholder && (inputPlaceholder.includes('需求') || inputPlaceholder.includes('想买'))) {
      console.log('   ✓ 输入框为自由需求输入');
    } else {
      console.log('   ⚠️  输入框可能不是自由需求输入');
    }
    
    // ========== 步骤3: 检查快捷示例是否为需求描述 ==========
    console.log('✅ 步骤3: 检查快捷示例是否为需求描述');
    const exampleButtons = await page.locator('button').all();
    let foundExamples = [];
    
    for (const button of exampleButtons) {
      const text = await button.textContent();
      if (text && (text.includes('想买') || text.includes('需要') || text.includes('键盘') || text.includes('耳机'))) {
        foundExamples.push(text.trim());
      }
    }
    
    if (foundExamples.length > 0) {
      console.log('   ✓ 找到快捷示例:');
      foundExamples.forEach(ex => console.log(`     - "${ex}"`));
    } else {
      console.log('   ⚠️  未找到明显的需求描述示例');
    }
    
    // ========== 步骤4: 输入测试需求 ==========
    console.log('✅ 步骤4: 输入测试需求 "想买个蓝牙耳机，通勤降噪"');
    const inputField = page.locator('input[type="text"], textarea').first();
    await inputField.fill('想买个蓝牙耳机，通勤降噪');
    await page.waitForTimeout(500);
    console.log('   ✓ 需求已输入');
    
    // ========== 步骤5: 点击"开始分析"按钮 ==========
    console.log('✅ 步骤5: 点击"开始分析"按钮（箭头图标按钮）');
    const analyzeButton = page.locator('button:has(svg)').first();
    await analyzeButton.click();
    console.log('   ✓ 已点击按钮');
    
    // ========== 步骤6: 验证跳转到确认页 ==========
    console.log('✅ 步骤6: 验证跳转到确认页 /confirm?requirement=...');
    await page.waitForURL(/\/confirm\?requirement=/, { timeout: 10000 });
    const currentURL = page.url();
    console.log(`   ✓ 已跳转到: ${currentURL}`);
    
    // ========== 步骤7: 等待 API 加载完成 ==========
    console.log('✅ 步骤7: 等待 API 加载完成');
    
    // 等待加载动画消失（最多等待60秒）
    try {
      await page.waitForSelector('.animate-spin', { state: 'detached', timeout: 60000 });
      console.log('   ✓ 加载动画已消失');
    } catch (e) {
      console.log('   ⚠️  加载动画未消失（超时60秒），继续检查内容');
    }
    
    // 等待关键内容出现
    try {
      await page.waitForSelector('text=/商品类型|品牌|评价维度/', { timeout: 10000 });
      console.log('   ✓ API 数据已加载');
    } catch (e) {
      console.log('   ⚠️  等待超时，但继续检查页面内容');
    }
    
    // 额外等待确保内容完全渲染
    await page.waitForTimeout(2000);
    
    // ========== 步骤8: 检查确认页显示内容 ==========
    console.log('✅ 步骤8: 检查确认页显示内容');
    
    const pageContent = await page.content();
    
    // 8.1 检查 AI 理解描述
    const hasAIUnderstanding = pageContent.includes('我理解您') || pageContent.includes('理解');
    console.log(`   ${hasAIUnderstanding ? '✓' : '✗'} AI 理解描述: ${hasAIUnderstanding ? '已显示' : '未找到'}`);
    
    // 8.2 检查商品类型
    const hasProductType = pageContent.includes('商品类型') || pageContent.includes('类型');
    console.log(`   ${hasProductType ? '✓' : '✗'} 商品类型: ${hasProductType ? '已显示' : '未找到'}`);
    
    // 8.3 检查品牌标签
    const hasBrands = pageContent.includes('品牌') || pageContent.includes('推荐品牌');
    console.log(`   ${hasBrands ? '✓' : '✗'} 品牌标签: ${hasBrands ? '已显示' : '未找到'}`);
    
    // 8.4 检查评价维度卡片
    const hasDimensions = pageContent.includes('评价维度') || pageContent.includes('维度');
    console.log(`   ${hasDimensions ? '✓' : '✗'} 评价维度卡片: ${hasDimensions ? '已显示' : '未找到'}`);
    
    // 8.5 检查搜索关键词
    const hasKeywords = pageContent.includes('搜索关键词') || pageContent.includes('关键词');
    console.log(`   ${hasKeywords ? '✓' : '✗'} 搜索关键词: ${hasKeywords ? '已显示' : '未找到'}`);
    
    // ========== 步骤9: 截图保存确认页面 ==========
    console.log('✅ 步骤9: 截图保存确认页面');
    await page.screenshot({ 
      path: 'screenshot-confirm.png',
      fullPage: true 
    });
    console.log('   ✓ 截图已保存: screenshot-confirm.png');
    
    // ========== 测试总结 ==========
    console.log('\n' + '='.repeat(50));
    console.log('🎉 测试完成！');
    console.log('='.repeat(50));
    console.log('测试结果汇总:');
    console.log(`  - 首页访问: ✓`);
    console.log(`  - 自由需求输入: ${inputPlaceholder ? '✓' : '?'}`);
    console.log(`  - 快捷示例: ${foundExamples.length > 0 ? '✓' : '?'}`);
    console.log(`  - 需求提交: ✓`);
    console.log(`  - 页面跳转: ✓`);
    console.log(`  - AI理解: ${hasAIUnderstanding ? '✓' : '✗'}`);
    console.log(`  - 商品类型: ${hasProductType ? '✓' : '✗'}`);
    console.log(`  - 品牌标签: ${hasBrands ? '✓' : '✗'}`);
    console.log(`  - 评价维度: ${hasDimensions ? '✓' : '✗'}`);
    console.log(`  - 搜索关键词: ${hasKeywords ? '✓' : '✗'}`);
    console.log(`  - 截图保存: ✓`);
    console.log('='.repeat(50));
    
  } catch (error) {
    console.error('\n❌ 测试过程中出现错误:');
    console.error(error.message);
    console.error('\n错误堆栈:');
    console.error(error.stack);
    
    // 出错时也截图
    try {
      await page.screenshot({ path: 'screenshot-error.png', fullPage: true });
      console.log('\n已保存错误截图: screenshot-error.png');
    } catch (e) {
      console.error('无法保存错误截图:', e.message);
    }
  } finally {
    // 关闭浏览器
    await browser.close();
    console.log('\n浏览器已关闭');
  }
})();
