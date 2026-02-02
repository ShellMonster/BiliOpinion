import React, { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import type { SentimentStats } from '../../../types/report';

interface SentimentPieProps {
  data: SentimentStats;
  title?: string;
}

/**
 * 情感分布饼图组件
 * 显示好评、中性评、差评的分布情况
 */
export const SentimentPie: React.FC<SentimentPieProps> = ({ data, title = '情感分布' }) => {
  const option = useMemo(() => {
    return {
      title: {
        text: title,
        left: 'center',
        textStyle: {
          fontSize: 16,
          fontWeight: 'normal'
        }
      },
      tooltip: {
        trigger: 'item',
        formatter: '{b}: {c} ({d}%)'
      },
      legend: {
        bottom: '0%',
        left: 'center'
      },
      series: [
        {
          name: '情感分布',
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: true,
            formatter: '{b}: {d}%'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 14,
              fontWeight: 'bold'
            }
          },
          labelLine: {
            show: true
          },
          data: [
            { 
              value: data.positive_count, 
              name: '好评',
              itemStyle: { color: '#22c55e' } // 绿色
            },
            { 
              value: data.neutral_count, 
              name: '中性',
              itemStyle: { color: '#9ca3af' } // 灰色
            },
            { 
              value: data.negative_count, 
              name: '差评',
              itemStyle: { color: '#ef4444' } // 红色
            }
          ]
        }
      ]
    };
  }, [data, title]);

  return (
    <div className="w-full h-[300px] bg-white rounded-lg p-4">
      <ReactECharts 
        option={option} 
        style={{ height: '100%', width: '100%' }}
        opts={{ renderer: 'svg' }}
      />
    </div>
  );
};
