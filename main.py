import requests
from bs4 import BeautifulSoup
import datetime
import math
import time


class Scrape:
    def __init__(self, tag):
        self.tag = tag
        self.url = f'https://stackoverflow.com/questions/tagged/{self.tag}?tab=frequent&page=1&pagesize=50'

    def request(self):
        response = requests.get(self.url)
        soup = BeautifulSoup(response.text, "html.parser")
        return soup

    # Title, views, votes, date of question, link

    def get_questions(self, soup):
        data = []
        try:
            qbody = soup.find("div", {"id": "questions"})
            questions = qbody.find_all(
                "div", {"class": "s-post-summary js-post-summary"})
            for question in questions:
                title = question.find(
                    "div", {"class": "s-post-summary--content"}).find("a").text
                # remove commas from title
                # title = title.replace(",", "\t")
                stats = question.find_all(
                    "div", {"class": "s-post-summary--stats-item"})
                votes = stats[0].find(
                    "span", {"class": "s-post-summary--stats-item-number"}).text
                answers = stats[1].find(
                    "span", {"class": "s-post-summary--stats-item-number"}).text
                views = stats[2].find(
                    "span", {"class": "s-post-summary--stats-item-number"}).text
                # get date from title attribute name
                # check if class relavetime exists if not date = 2010-12-1 00:00:00Z
                if question.find("span", {"class": "relativetime"}) is None:
                    date = "2010-12-1 00:00:00Z"
                else:
                    date = question.find(
                        "span", {"class": "relativetime"})["title"]
                # convert date (2014-12-16 09:00:51Z) to datetime object
                date = datetime.datetime.strptime(date, "%Y-%m-%d %H:%M:%S%z")
                # date to unix timestamp
                date = date.timestamp()

                link = question.find(
                    "div", {"class": "s-post-summary--content"}).find("a")["href"]
                link = "https://stackoverflow.com" + link

                # Calculate views per day given the date it was posted todate and views
                # if views has k, multiply by 1000
                if "k" in views:
                    views = int(float(views.replace("k", "")) * 1000)
                elif "m" in views:
                    views = int(float(views.replace("m", "")) * 1000000)
                else:
                    views = int(views)
                # get current date
                now = datetime.datetime.now()
                # convert to unix timestamp
                now = now.timestamp()
                # calculate days since posted
                days = (now - date) / 86400
                # calculate views per day
                views_per_day = math.floor(views / days)

                data.append(
                    [title, votes, answers, views, views_per_day, link])
        except:
            pass
        return data

    def get_AllPages(self, soup):
        # get total number of pages
        data = []

        pages = soup.find(
            "div", {"class": "s-pagination site1 themed pager float-left"}).find_all("a")[-2].text
        # loop through all pages
        pages = 45
        for page in range(1, int(pages) + 1):
            # get url
            url = f'https://stackoverflow.com/questions/tagged/{self.tag}?tab=frequent&page={page}&pagesize=50'
            # get response
            response = requests.get(url)
            # get soup
            soup = BeautifulSoup(response.text, "html.parser")
            # get questions
            data += self.get_questions(soup)

        return data

    def scrape(self):
        soup = self.request()
        data = self.get_AllPages(soup)

        try:
            with open(f"{self.tag}_data.csv", "w") as f:
                f.write("Title, Votes, Answers, Views, Views per day, Link \n")
                for row in data:
                    # for the title include quotes
                    f.write(
                        f'"{row[0]}", {row[1]}, {row[2]}, {row[3]}, {row[4]}, {row[5]} \n')
        except:
            pass


def main():
    tag = input("Enter tag: ")
    starttime = time.time()
    scraper = Scrape(tag)
    scraper.scrape()
    print(f"Time taken Python: {time.time() - starttime}")


if __name__ == "__main__":
    main()
