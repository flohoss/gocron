document.addEventListener('alpine:init', () => {
    Alpine.data('event', () => ({
        runIcon: '',
        idle: false,
        SSE: null,
        data: null,
        init() {
            this.SSE = new EventSource("/api/events?stream=status");
            this.SSE.onmessage = (event) => {
                this.data = JSON.parse(event.data);
                console.log(this.data);
                this.handleData();
            };
            window.addEventListener('beforeunload', () => {
                this.$nextTick(() => this.cleanup());
            });
        },
        cleanup() {
            if (this.SSE) {
                this.SSE.close();
            }
        },
        handleData() {
            this.idle = this.data.idle;
            this.handleRunIcon();
        },
        handleRunIcon() {
            if (this.idle) {
                this.runIcon = `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
                                    <path stroke-linecap="round" stroke-linejoin="round" d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.347a1.125 1.125 0 0 1 0 1.972l-11.54 6.347a1.125 1.125 0 0 1-1.667-.986V5.653Z"></path>
                                </svg>`;
            } else {
                this.runIcon = `<span class="loading loading-spinner"></span>`;
            }
        },
        run() {
            fetch('/api/jobs', { method: 'POST' })
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`HTTP error! status: ${response.status}`);
                    }
                });
        }
    }));
});