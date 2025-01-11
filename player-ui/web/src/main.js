import express from 'express';
import path from 'path';
import cors from 'cors';
import fetch from 'node-fetch';

const app = express();
const PORT = 3000;

app.use(cors());
// Serve static files from 'public' and 'dist'
app.use(express.static(path.join(process.cwd(), 'public')));
app.use(express.static(path.join(process.cwd(), 'dist')));

// Default route to serve index.html
app.get('/', (req, res) => {
    res.sendFile(path.join(process.cwd(), 'public', 'index.html'));
});
app.get('/proxy-image', async (req, res) => {
    const imageUrl = req.query.url;

    if (!imageUrl) {
        return res.status(400).send('Image URL is required');
    }

    try {
        const response = await fetch(imageUrl);
        if (!response.ok) {
            return res.status(response.status).send('Error fetching image');
        }

        res.set('Content-Type', response.headers.get('content-type'));
        response.body.pipe(res);
    } catch (error) {
        console.error('Error fetching image:', error);
        res.status(500).send('Internal server error');
    }
});
// Start the server
app.listen(PORT, () => {
    console.log(`Server is running at http://localhost:${PORT}`);
});


