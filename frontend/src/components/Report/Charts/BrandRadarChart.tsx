import React, { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import type { ReportData } from '../../../types/report';

interface BrandRadarChartProps {
  data: ReportData;
}

/**
 * 品牌多维度对比雷达图组件
 * @param data 报告数据
 */
export const BrandRadarChart: React.FC<BrandRadarChartProps> = ({ data }) => {
  // 品牌颜色配置
  const colors = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981'];

  // 使用 useMemo 优化图表配置计算
  const option = useMemo(() => ({
    tooltip: {
      trigger: 'item'
    },
    legend: {
      data: data.brands.slice(0, 3),
      bottom: 10,
      textStyle: {
        fontSize: 12
      }
    },
    radar: {
      indicator: data.dimensions.map(dim => ({
        name: dim.name,
        max: 100
      })),
      splitNumber: 4,
      name: {
        textStyle: {
          color: '#6b7280',
          fontSize: 11
        }
      },
      splitLine: {
        lineStyle: {
          color: '#e5e7eb'
        }
      },
      splitArea: {
        show: true,
        areaStyle: {
          color: ['rgba(255, 255, 255, 0)', 'rgba(249, 250, 251, 0.5)']
        }
      },
      axisLine: {
        lineStyle: {
          color: '#e5e7eb'
        }
      }
    },
    series: [{
      type: 'radar',
      data: data.brands.slice(0, 3).map((brand, index) => ({
        value: data.dimensions.map(dim => 
          data.scores[brand]?.[dim.name] ? data.scores[brand][dim.name] * 10 : 0
        ),
        name: brand,
        lineStyle: {
          color: colors[index],
          width: 2
        },
        areaStyle: {
          color: colors[index],
          opacity: 0.2
        },
        itemStyle: {
          color: colors[index]
        }
      }))
    }]
  }), [data, colors]);

  return <ReactECharts option={option} style={{ height: '320px' }} />;
};
