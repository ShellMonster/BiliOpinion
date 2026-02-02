import React, { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import type { ReportData } from '../../../types/report';

interface BrandScoreChartProps {
  data: ReportData;
}

/**
 * 品牌综合得分排名柱状图组件
 * @param data 报告数据
 */
export const BrandScoreChart: React.FC<BrandScoreChartProps> = ({ data }) => {
  // 品牌颜色配置
  const colors = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981'];

  // 使用 useMemo 优化图表配置计算
  const option = useMemo(() => ({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: 80,
      right: 30,
      top: 20,
      bottom: 20
    },
    xAxis: {
      type: 'value',
      max: 10,
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: '#6b7280'
      },
      splitLine: {
        lineStyle: {
          color: '#f0f0f0'
        }
      }
    },
    yAxis: {
      type: 'category',
      data: [...data.rankings].reverse().map(r => r.brand),
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: '#6b7280'
      }
    },
    series: [{
      type: 'bar',
      data: [...data.rankings].reverse().map((r, index) => ({
        value: Math.round(r.overall_score * 10) / 10,
        itemStyle: {
          color: colors[index % colors.length],
          borderRadius: [0, 8, 8, 0]
        }
      })),
      barWidth: 30,
      label: {
        show: true,
        position: 'right',
        formatter: '{c}',
        color: '#374151',
        fontWeight: 'bold'
      }
    }]
  }), [data, colors]);

  return <ReactECharts option={option} style={{ height: '320px' }} />;
};
