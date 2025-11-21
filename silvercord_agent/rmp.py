from ratemyprofessor import Professor
import requests
import bs4 as bs
import re
from urllib.parse import quote_plus
from constants import CCNY_SCHOOL_ID

headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
}


def query_rmp(college_id: str, professor_name: str) -> (list, bool):
    encoded_name = quote_plus(professor_name)
    url = 'https://www.ratemyprofessors.com/search/professors/%s?q=%s' % (college_id, encoded_name)
    response = requests.get(url, headers=headers)
    soup = bs.BeautifulSoup(response.text, 'html.parser')
    
    prof_links = soup.find_all('a', href=re.compile(r'/professor/\d+'))
    professor_ids = []
    for link in prof_links:
        href = link.get('href', '')
        match = re.search(r'/professor/(\d+)', href)
        if match:
            prof_id = match.group(1)
            if prof_id != college_id and prof_id not in professor_ids:  # Filter out school ID and duplicates
                professor_ids.append(prof_id)
    
    return professor_ids, len(professor_ids) > 0


def get_professor_url(professor_name: str) -> str | None:
    professor_ids, found = query_rmp(CCNY_SCHOOL_ID, professor_name)
    if found and professor_ids:
        return f"https://www.ratemyprofessors.com/professor/{professor_ids[0]}"
    else:
        return None