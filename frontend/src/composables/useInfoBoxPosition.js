import { ref, watch } from 'vue';

export function useInfoBoxPosition(mapViewComponentRef, selectedIsland) {
    const infoBoxStyle = ref({ display: 'none' });
    const transformCounter = ref(0);

    const calculateInfoBoxStyle = () => {
        const svgElement = mapViewComponentRef.value?.svgRef;
        if (!selectedIsland.value || !svgElement) {
            infoBoxStyle.value = { display: 'none' };
            return;
        }

        const island = selectedIsland.value;
        const pt = svgElement.createSVGPoint();
        pt.x = island.x;
        pt.y = island.y;
        const screenPoint = pt.matrixTransform(svgElement.getScreenCTM());

        infoBoxStyle.value = {
            position: 'fixed',
            top: `${screenPoint.y}px`,
            left: `${screenPoint.x}px`,
            transform: 'translate(-50%, -100%) translateY(-20px)',
        };
    };

    watch([selectedIsland, transformCounter], calculateInfoBoxStyle, {
        flush: 'post',
    });

    const updateInfoBoxPosition = () => {
        transformCounter.value++;
    };

    return { infoBoxStyle, updateInfoBoxPosition };
}
