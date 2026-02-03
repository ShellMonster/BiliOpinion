import React, { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import type { SentimentStats } from '../../../types/report';

interface ScoreHistogramProps {
  data: Record<string, SentimentStats>;
  title?: string;
}

interface EChartsLabelParams {
  value: number
  dataIndex: number
  seriesIndex: number
}

/**
 * 品牌评分分布柱状图组件
 * 显示各品牌的好评、中评、差评数量分布
 */
export const ScoreHistogram: React.FC<ScoreHistogramProps> = ({ data, title = '品牌评分分布' }) => {
  const option = useMemo(() => {
    // 获取所有品牌名称
    const brands = Object.keys(data);
    
    // 提取各情感维度的数据
    const positiveCounts = brands.map(brand => data[brand].positive_count);
    const neutralCounts = brands.map(brand => data[brand].neutral_count);
    const negativeCounts = brands.map(brand => data[brand].negative_count);

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
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      legend: {
        data: ['好评', '中评', '差评'],
        bottom: '0%',
        left: 'center'
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '10%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: brands,
        axisTick: {
          alignWithLabel: true
        },
        axisLabel: {
          interval: 0,
          rotate: brands.length > 5 ? 30 : 0 // 品牌较多时旋转标签
        }
      },
      yAxis: {
        type: 'value',
        name: '评论数量'
      },
      series: [
        {
          name: '好评',
          type: 'bar',
          data: positiveCounts,
          itemStyle: {
            color: '#22c55e' // 绿色
          },
          label: {
            show: true,
            position: 'top',
            formatter: (params: EChartsLabelParams) => params.value > 0 ? params.value : ''
          }
        },
        {
          name: '中评',
          type: 'bar',
          data: neutralCounts,
          itemStyle: {
            color: '#9ca3af' // 灰色
          },
          label: {
            show: true,
            position: 'top',
            formatter: (params: EChartsLabelParams) => params.value > 0 ? params.value : ''
          }
        },
        {
          name: '差评',
          type: 'bar',
          data: negativeCounts,
          itemStyle: {
            color: '#ef4444' // 红色
          },
          label: {
            show: true,
            position: 'top',
            formatter: (params: EChartsLabelParams) => params.value > 0 ? params.value : ''
          }
        }
      ]
    };
  }, [data, title]);

  return (
    <div className="w-full h-[400px] bg-white rounded-lg p-4">
      <ReactECharts 
        option={option} 
        style={{ height: '100%', width: '100%' }}
        opts={{ renderer: 'svg' }}
      />
    </div>
  );
};
