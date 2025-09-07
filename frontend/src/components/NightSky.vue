<template>
    <canvas id="starCanvas" ref="canvas" class="fixed inset-0 z-0"></canvas>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref } from 'vue';

const canvas = ref(null);
let animationId;

onMounted(() => {
    const cvs = canvas.value;
    const ctx = cvs.getContext('2d');

    let mouseX = window.innerWidth / 2;
    let mouseY = window.innerHeight / 2;
    let parallaxX = 0;
    let parallaxY = 0;

    const NUM_STARS = 600;
    const MILKY_WAY_STAR_COUNT = 500;
    const SHOOTING_STAR_FREQ = Math.random() < 0.5 ? 0.01 : 0.03;
    const HOVER_RADIUS = 150;
    const MAX_FADE_OPACITY = 0.2;

    const milkyWay = {
        centerX: 0,
        centerY: 0,
        width: 0,
        length: 0,
        angle: -Math.PI / 20,
    };
    let stars = [],
        milkyWayStars = [],
        shootingStars = [];

    const clamp = (v, min, max) => Math.min(max, Math.max(min, v));

    function generateStarColor() {
        const t = Math.random();
        let r, g, b;
        if (t < 0.01) {
            r = 155 + Math.random() * 20;
            g = 176 + Math.random() * 20;
            b = 255;
        } else if (t < 0.1) {
            r = 170 + Math.random() * 30;
            g = 200 + Math.random() * 40;
            b = 255;
        } else if (t < 0.3) {
            r = 235 + Math.random() * 15;
            g = 240 + Math.random() * 15;
            b = 255;
        } else if (t < 0.5) {
            r = 255;
            g = 245 + Math.random() * 10;
            b = 220 + Math.random() * 20;
        } else if (t < 0.7) {
            r = 255;
            g = 230 + Math.random() * 10;
            b = 180 + Math.random() * 20;
        } else if (t < 0.9) {
            r = 255;
            g = 190 + Math.random() * 20;
            b = 120 + Math.random() * 20;
        } else {
            r = 255;
            g = 120 + Math.random() * 20;
            b = 100 + Math.random() * 20;
        }
        r = clamp(r + Math.random() * 8 - 4, 0, 255);
        g = clamp(g + Math.random() * 8 - 4, 0, 255);
        b = clamp(b + Math.random() * 8 - 4, 0, 255);
        return `${Math.round(r)}, ${Math.round(g)}, ${Math.round(b)}`;
    }

    function createStar() {
        const radius = Math.random() * 1.1 + 0.1;
        return {
            x: Math.random() * cvs.width,
            y: Math.random() * cvs.height,
            radius,
            baseOpacity: Math.random() * 0.7 + 0.3,
            opacity: 0,
            flickerSpeed: Math.random() * 0.1 + 0.02,
            flickerPhase: Math.random() * Math.PI * 2,
            vx: (Math.random() - 0.5) * 0.001,
            vy: (Math.random() - 0.5) * 0.001,
            color: generateStarColor(),
            glow: Math.random() < 0.25,
        };
    }

    function createMilkyWayStar() {
        const posX = (Math.random() - 0.5) * milkyWay.length;
        const curveY = 0.00002 * posX * posX;
        const posY = curveY + (Math.random() - 0.5) * milkyWay.width;
        const rotatedX =
            posX * Math.cos(milkyWay.angle) - posY * Math.sin(milkyWay.angle);
        const rotatedY =
            posX * Math.sin(milkyWay.angle) + posY * Math.cos(milkyWay.angle);
        const radius = Math.random() * 0.6 + 0.1;
        return {
            x: milkyWay.centerX + rotatedX,
            y: milkyWay.centerY + rotatedY,
            radius,
            baseOpacity: Math.random() * 0.6 + 0.4,
            opacity: 0,
            flickerSpeed: Math.random() * 0.1 + 0.05,
            flickerPhase: Math.random() * Math.PI * 2,
            vx: 0,
            vy: 0,
            color: generateStarColor(),
            glow: Math.random() < 0.3,
        };
    }

    function createShootingStar() {
        const speed = 30 + Math.random() * 30;
        const maxLife = 40 + Math.random() * 20;
        const length = speed * (6 + Math.random() * 6);

        const startX = Math.random() * cvs.width;
        const startY = Math.random() * cvs.height;

        const shootLeft = startX > cvs.width / 2;
        const baseAngle = shootLeft ? Math.PI : 0;
        const angleVariation = (Math.random() - 0.5) * (Math.PI / 10);
        const angle = baseAngle + angleVariation;

        const colors = [
            '180,230,255',
            '100,210,220',
            '80,200,180',
            '120,255,210',
            '140,255,240',
        ];
        const color = colors[Math.floor(Math.random() * colors.length)];

        return {
            x: startX,
            y: startY,
            speed,
            length,
            angle,
            opacity: 0.2,
            life: 0,
            maxLife,
            trail: [],
            color,
        };
    }

    function resizeCanvas() {
        cvs.width = window.innerWidth;
        cvs.height = window.innerHeight;
        milkyWay.centerX = cvs.width / 2;
        milkyWay.centerY = cvs.height / 2 + 50;
        milkyWay.width = cvs.height / 4;
        milkyWay.length = cvs.width * 1.2;
        stars = [];
        milkyWayStars = [];
        shootingStars = [];
        for (let i = 0; i < NUM_STARS; i++) stars.push(createStar());
        for (let i = 0; i < MILKY_WAY_STAR_COUNT; i++)
            milkyWayStars.push(createMilkyWayStar());
    }

    window.addEventListener('resize', resizeCanvas);
    resizeCanvas();

    window.addEventListener('mousemove', e => {
        mouseX = (e.clientX / cvs.width) * 2 - 1;
        mouseY = (e.clientY / cvs.height) * 2 - 1;
    });

    function drawStar(star, ox = 0, oy = 0) {
        const x = star.x + ox,
            y = star.y + oy;
        if (star.glow) {
            const glowR = star.radius * 5 + Math.random() * 3;
            const grad = ctx.createRadialGradient(
                x,
                y,
                star.radius,
                x,
                y,
                glowR
            );
            grad.addColorStop(0, `rgba(${star.color},${star.opacity * 0.6})`);
            grad.addColorStop(1, 'rgba(0,0,0,0)');
            ctx.fillStyle = grad;
            ctx.beginPath();
            ctx.arc(x, y, glowR, 0, Math.PI * 2);
            ctx.fill();
        }
        ctx.beginPath();
        ctx.fillStyle = `rgba(${star.color},${star.opacity})`;
        ctx.shadowColor = `rgba(${star.color},${star.opacity})`;
        ctx.shadowBlur = star.radius * 1.5 + (star.glow ? 1.5 : 0);
        ctx.arc(x, y, star.radius, 0, Math.PI * 2);
        ctx.fill();
        ctx.shadowBlur = 0;
    }

    function drawMilkyWayGlow() {
        const g = ctx.createRadialGradient(
            milkyWay.centerX,
            milkyWay.centerY,
            milkyWay.width / 4,
            milkyWay.centerX,
            milkyWay.centerY,
            milkyWay.width
        );
        g.addColorStop(0, 'rgba(255,255,255,0.25)');
        g.addColorStop(1, 'rgba(255,255,255,0)');
        ctx.fillStyle = g;
        ctx.save();
        ctx.translate(milkyWay.centerX, milkyWay.centerY);
        ctx.rotate(milkyWay.angle);
        ctx.beginPath();
        ctx.ellipse(
            0,
            0,
            milkyWay.length / 2,
            milkyWay.width,
            0,
            0,
            Math.PI * 2
        );
        ctx.fill();
        ctx.restore();
    }

    function getFadeOpacity(x, y) {
        const mx = ((mouseX + 1) / 2) * cvs.width,
            my = ((mouseY + 1) / 2) * cvs.height;
        const dx = x - mx,
            dy = y - my,
            dist = Math.sqrt(dx * dx + dy * dy);
        if (dist < HOVER_RADIUS) {
            return (
                MAX_FADE_OPACITY +
                (dist / HOVER_RADIUS) * (1 - MAX_FADE_OPACITY)
            );
        }
        return 1;
    }

    function drawShootingStar(shoot) {
        ctx.save();
        ctx.translate(shoot.x, shoot.y);
        ctx.rotate(shoot.angle);

        const tailLength = shoot.length;
        const gradient = ctx.createLinearGradient(0, 0, -tailLength, 0);
        gradient.addColorStop(0, `rgba(${shoot.color},${shoot.opacity * 0.7})`);
        gradient.addColorStop(
            0.7,
            `rgba(${shoot.color},${shoot.opacity * 0.4})`
        );
        gradient.addColorStop(1, `rgba(${shoot.color},0)`);

        ctx.fillStyle = gradient;
        ctx.beginPath();
        ctx.moveTo(0, -1.2);
        ctx.lineTo(-tailLength, -0.15);
        ctx.lineTo(-tailLength, 0.15);
        ctx.lineTo(0, 1.2);
        ctx.closePath();
        ctx.fill();

        ctx.shadowColor = `rgba(255,255,255,${shoot.opacity})`;
        ctx.shadowBlur = 25;
        ctx.beginPath();
        ctx.fillStyle = `rgba(255,255,255,${shoot.opacity})`;
        ctx.arc(0, 0, 3, 0, Math.PI * 2);
        ctx.fill();
        ctx.restore();

        for (let i = shoot.trail.length - 1; i >= 0; i--) {
            const p = shoot.trail[i];
            p.opacity -= 0.03;
            p.radius *= 0.9;
            if (p.opacity <= 0 || p.radius <= 0.1) {
                shoot.trail.splice(i, 1);
                continue;
            }
            ctx.beginPath();
            ctx.fillStyle = `rgba(${shoot.color},${p.opacity})`;
            ctx.arc(p.x, p.y, p.radius, 0, Math.PI * 2);
            ctx.fill();
        }
    }

    function animate() {
        ctx.clearRect(0, 0, cvs.width, cvs.height);
        parallaxX += (mouseX * 50 - parallaxX) * 0.05;
        parallaxY += (mouseY * 30 - parallaxY) * 0.05;
        drawMilkyWayGlow();

        const drift = 0.00005;
        for (const star of stars) {
            star.flickerPhase += star.flickerSpeed;
            const base = clamp(
                star.baseOpacity + Math.sin(star.flickerPhase) * 0.15,
                0.3,
                1
            );
            const ox = parallaxX * star.radius * 0.3,
                oy = parallaxY * star.radius * 0.3;
            star.opacity = base * getFadeOpacity(star.x + ox, star.y + oy);
            star.x += star.vx + drift;
            star.y += star.vy;
            if (star.x > cvs.width) star.x = 0;
            else if (star.x < 0) star.x = cvs.width;
            if (star.y > cvs.height) star.y = 0;
            else if (star.y < 0) star.y = cvs.height;
            drawStar(star, ox, oy);
        }

        for (const star of milkyWayStars) {
            star.flickerPhase += star.flickerSpeed;
            const base = clamp(
                star.baseOpacity + Math.sin(star.flickerPhase) * 0.15,
                0.3,
                1
            );
            const ox = parallaxX * star.radius * 0.2,
                oy = parallaxY * star.radius * 0.2;
            star.opacity = base * getFadeOpacity(star.x + ox, star.y + oy);
            drawStar(star, ox, oy);
        }

        // shooting stars
        if (shootingStars.length === 0 && Math.random() < SHOOTING_STAR_FREQ) {
            shootingStars.push(createShootingStar());
        }
        for (let i = shootingStars.length - 1; i >= 0; i--) {
            const s = shootingStars[i];
            s.life++;
            s.trail.push({
                x: s.x,
                y: s.y,
                radius: 0.6 + Math.random() * 0.3,
                opacity: s.opacity,
            });
            s.x += Math.cos(s.angle) * s.speed;
            s.y += Math.sin(s.angle) * s.speed;
            if (s.life < 30) s.opacity = (s.life / 30) ** 2;
            else if (s.life > s.maxLife - 10)
                s.opacity = (s.maxLife - s.life) / 10;
            else s.opacity = 1;
            drawShootingStar(s);
            if (s.opacity <= 0 || s.life > s.maxLife)
                shootingStars.splice(i, 1);
        }

        animationId = requestAnimationFrame(animate);
    }
    animate();
});

onBeforeUnmount(() => {
    cancelAnimationFrame(animationId);
});
</script>

<style scoped>
#starCanvas {
    background:
        radial-gradient(
            ellipse 60% 20% at 50% 40%,
            rgba(255, 255, 255, 0.01),
            transparent 80%
        ),
        radial-gradient(
            ellipse at center,
            #0a0f1c 0%,
            #050812 60%,
            #000000 100%
        );
    background-blend-mode: screen;
}
</style>
