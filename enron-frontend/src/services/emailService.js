import axios from 'axios';
import { API_CONFIG } from '../constants/config';

export const searchEmails = async ({ query, from, to, page, size }) => {
  try {
    const payload = {
      query: query || '',
      from: page * size,
      size: size
    };

    if (from || to) {
      payload.filters = {};
      if (from) payload.filters.from = `*${from}*`;
      if (to) payload.filters.to = `*${to}*`;
    }

    const response = await axios.post(
      `${API_CONFIG.BASE_URL}${API_CONFIG.ENDPOINTS.SEARCH}`, 
      payload
    );
    return response.data;
  } catch (error) {
    console.error("Error searching emails:", error);
    throw error;
  }
};
