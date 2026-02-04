import React, { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import 'echarts-wordcloud';
import type { KeywordItem } from '../../../types/report';

interface KeywordCloudProps {
  data: KeywordItem[];
}

export const KeywordCloud: React.FC<KeywordCloudProps> = ({ data }) => {
  const chartOption = useMemo(() => {
    // 转换数据格式：{ name: word, value: count }
    // 既然要求"大小与频率成正比"，echarts-wordcloud 默认即是如此
    // 限制显示前50个高频关键词
    const formattedData = [...data]
      .sort((a, b) => b.count - a.count)
      .slice(0, 50)
      .map((item) => ({
        name: item.word,
        value: item.count,
      }));

    return {
      tooltip: {
        show: true,
        formatter: '{b}: {c}次', // 显示词频
      },
      series: [
        {
          type: 'wordCloud',
          // 形状：'circle', 'cardioid', 'diamond', 'triangle-forward', 'triangle', 'pentagon', 'star'
          shape: 'circle',
          
          // 布局位置和大小
          left: 'center',
          top: 'center',
          width: '95%',
          height: '95%',
          right: null,
          bottom: null,

          // 词的大小范围
          sizeRange: [14, 60],

          // 词的旋转范围和步长
          rotationRange: [-45, 90],
          rotationStep: 45,

          // 词之间的间距
          gridSize: 8,

          // 是否允许词超出画布范围
          drawOutOfBound: false,

          // 布局动画
          layoutAnimation: true,

          // 全局文本样式
          textStyle: {
            fontFamily: 'sans-serif',
            fontWeight: 'bold',
            // 随机颜色
            color: function () {
              return 'rgb(' + [
                Math.round(Math.random() * 160),
                Math.round(Math.random() * 160),
                Math.round(Math.random() * 160)
              ].join(',') + ')';
            }
          },
          
          // 高亮样式
          emphasis: {
            focus: 'self',
            textStyle: {
              shadowBlur: 10,
              shadowColor: '#333'
            }
          },

          // 数据
          data: formattedData,
        },
      ],
    };
  }, [data]);

  if (!data || data.length === 0) {
    return <div className="flex h-64 items-center justify-center text-gray-400">暂无关键词数据</div>;
  }

  return (
    <div className="w-full bg-white rounded-lg p-4 shadow-sm border border-gray-200">
      <h3 className="text-lg font-semibold text-gray-800 mb-4">关键词云</h3>
      <div className="h-[400px]">
        <ReactECharts
          option={chartOption}
          style={{ height: '100%', width: '100%' }}
          notMerge={true}
          lazyUpdate={true}
        />
      </div>
    </div>
  );
};
